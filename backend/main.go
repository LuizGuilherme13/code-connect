// @title           Code Connect API
// @version         1.0
// @description     API do projeto Code Connect
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}" (e.g. Bearer eyJhbGci...)
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"code-connect/backend/handlers"
	"code-connect/backend/middleware"
	"code-connect/backend/store"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "code-connect/backend/docs"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	if err := store.InitSchema(db); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	s := store.New(db)
	h := handlers.New(s)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.Handle("GET /api/users/me", middleware.AuthMiddleware(http.HandlerFunc(h.GetProfile)))

	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/swagger.json")
	})
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	handler := middleware.CorsMiddleware(mux)

	log.Printf("Backend running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
