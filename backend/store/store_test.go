package store

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

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

	if err := InitSchema(db); err != nil {
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

func TestCreateUser(t *testing.T) {
	db := connectDB(t)
	defer db.Close()
	truncateUsers(t, db)

	s := New(db)

	user, err := s.Create("João", "joao@email.com", "hashed_password")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == "" {
		t.Error("expected user ID to be set")
	}
	if user.Name != "João" {
		t.Errorf("expected name 'João', got '%s'", user.Name)
	}
	if user.Email != "joao@email.com" {
		t.Errorf("expected email 'joao@email.com', got '%s'", user.Email)
	}
	if user.Password != "hashed_password" {
		t.Errorf("expected password 'hashed_password', got '%s'", user.Password)
	}
}

func TestFindByEmail(t *testing.T) {
	db := connectDB(t)
	defer db.Close()
	truncateUsers(t, db)

	s := New(db)

	s.Create("João", "joao@email.com", "hashed_password")

	t.Run("existing email", func(t *testing.T) {
		user, err := s.FindByEmail("joao@email.com")
		if err != nil {
			t.Fatalf("FindByEmail failed: %v", err)
		}
		if user.Name != "João" {
			t.Errorf("expected name 'João', got '%s'", user.Name)
		}
	})

	t.Run("non-existing email", func(t *testing.T) {
		_, err := s.FindByEmail("naoexiste@email.com")
		if err != ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestFindByID(t *testing.T) {
	db := connectDB(t)
	defer db.Close()
	truncateUsers(t, db)

	s := New(db)

	created, err := s.Create("João", "joao@email.com", "hashed_password")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	t.Run("existing ID", func(t *testing.T) {
		user, err := s.FindByID(created.ID)
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}
		if user.Email != "joao@email.com" {
			t.Errorf("expected email 'joao@email.com', got '%s'", user.Email)
		}
	})

	t.Run("non-existing ID", func(t *testing.T) {
		_, err := s.FindByID("999")
		if err != ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	db := connectDB(t)
	defer db.Close()
	truncateUsers(t, db)

	s := New(db)

	_, err := s.Create("João", "joao@email.com", "pass")
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = s.Create("João", "joao@email.com", "pass")
	if err != ErrEmailExists {
		t.Errorf("expected ErrEmailExists, got %v", err)
	}
}

func TestCreateUserIncrementalID(t *testing.T) {
	db := connectDB(t)
	defer db.Close()
	truncateUsers(t, db)

	s := New(db)

	u1, err := s.Create("User1", "u1@email.com", "pass")
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	u2, err := s.Create("User2", "u2@email.com", "pass")
	if err != nil {
		t.Fatalf("second create failed: %v", err)
	}

	if u1.ID == u2.ID {
		t.Errorf("expected different IDs, got '%s' for both", u1.ID)
	}
}

func TestFindByEmailAfterMultipleCreates(t *testing.T) {
	db := connectDB(t)
	defer db.Close()
	truncateUsers(t, db)

	s := New(db)

	s.Create("João", "joao@email.com", "pass1")
	s.Create("Maria", "maria@email.com", "pass2")

	user, err := s.FindByEmail("maria@email.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if user.Name != "Maria" {
		t.Errorf("expected name 'Maria', got '%s'", user.Name)
	}
}
