package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"code-connect/backend/middleware"
	"code-connect/backend/models"
	"code-connect/backend/store"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

const testSecret = "test-secret-key"

func connectDB(t *testing.T) *sql.DB {
	t.Helper()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "codeconnect")
	dbPass := getEnv("DB_PASSWORD", "codeconnect")
	dbName := getEnv("DB_NAME", "codeconnect")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	if err := store.InitSchema(db); err != nil {
		t.Fatalf("failed to init schema: %v", err)
	}

	return db
}

func truncateUsers(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("failed to truncate users: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func setupTestHandler(t *testing.T) (*Handler, *sql.DB) {
	t.Helper()
	db := connectDB(t)
	truncateUsers(t, db)
	s := store.New(db)
	return New(s), db
}

func performRequest(h *Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.Handle("GET /api/users/me", middleware.AuthMiddleware(http.HandlerFunc(h.GetProfile)))

	mux.ServeHTTP(rr, req)
	return rr
}

func performRequestWithAuth(h *Handler, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("GET /api/users/me", middleware.AuthMiddleware(http.HandlerFunc(h.GetProfile)))

	mux.ServeHTTP(rr, req)
	return rr
}

func TestRegister_Success(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	body := models.RegisterRequest{
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}

	rr := performRequest(h, "POST", "/api/register", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var resp models.UserResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Name != "João" {
		t.Errorf("expected name 'João', got '%s'", resp.Name)
	}
	if resp.Email != "joao@email.com" {
		t.Errorf("expected email 'joao@email.com', got '%s'", resp.Email)
	}
	if resp.ID == "" {
		t.Error("expected user ID to be set")
	}
}

func TestRegister_MissingFields(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	tests := []struct {
		name string
		body models.RegisterRequest
	}{
		{"missing name", models.RegisterRequest{Email: "a@b.com", Password: "123"}},
		{"missing email", models.RegisterRequest{Name: "João", Password: "123"}},
		{"missing password", models.RegisterRequest{Name: "João", Email: "a@b.com"}},
		{"all empty", models.RegisterRequest{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := performRequest(h, "POST", "/api/register", tt.body)
			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rr.Code)
			}
		})
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	body := models.RegisterRequest{
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}

	rr1 := performRequest(h, "POST", "/api/register", body)
	if rr1.Code != http.StatusCreated {
		t.Fatalf("first register failed: %d", rr1.Code)
	}

	rr2 := performRequest(h, "POST", "/api/register", body)
	if rr2.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", rr2.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	registerBody := models.RegisterRequest{
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}
	performRequest(h, "POST", "/api/register", registerBody)

	loginBody := models.LoginRequest{
		Email:    "joao@email.com",
		Password: "123456",
	}
	rr := performRequest(h, "POST", "/api/login", loginBody)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp models.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Token == "" {
		t.Error("expected token to be set")
	}

	token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil || !token.Valid {
		t.Errorf("token is not valid: %v", err)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	loginBody := models.LoginRequest{
		Email:    "naoexiste@email.com",
		Password: "123456",
	}
	rr := performRequest(h, "POST", "/api/login", loginBody)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	registerBody := models.RegisterRequest{
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}
	performRequest(h, "POST", "/api/register", registerBody)

	loginBody := models.LoginRequest{
		Email:    "joao@email.com",
		Password: "wrongpassword",
	}
	rr := performRequest(h, "POST", "/api/login", loginBody)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetProfile_Success(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	registerBody := models.RegisterRequest{
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	}
	rr := performRequest(h, "POST", "/api/register", registerBody)

	var registered models.UserResponse
	json.NewDecoder(rr.Body).Decode(&registered)

	token := generateTestToken(t, registered.ID, testSecret, time.Hour)

	rr2 := performRequestWithAuth(h, "GET", "/api/users/me", token)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr2.Code)
	}

	var profile models.UserResponse
	if err := json.NewDecoder(rr2.Body).Decode(&profile); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if profile.Name != "João" {
		t.Errorf("expected name 'João', got '%s'", profile.Name)
	}
	if profile.Email != "joao@email.com" {
		t.Errorf("expected email 'joao@email.com', got '%s'", profile.Email)
	}
}

func TestGetProfile_NoToken(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/users/me", nil)
	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("GET /api/users/me", middleware.AuthMiddleware(http.HandlerFunc(h.GetProfile)))
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetProfile_InvalidToken(t *testing.T) {
	h, db := setupTestHandler(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("GET /api/users/me", middleware.AuthMiddleware(http.HandlerFunc(h.GetProfile)))
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func generateTestToken(t *testing.T, userID, secret string, expiration time.Duration) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(expiration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return tokenString
}
