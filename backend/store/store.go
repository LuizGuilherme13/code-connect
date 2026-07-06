package store

import (
	"database/sql"
	"errors"
	"fmt"

	"code-connect/backend/models"

	"github.com/lib/pq"
)

var (
	ErrEmailExists       = errors.New("email already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type PGStore struct {
	db *sql.DB
}

func New(db *sql.DB) *PGStore {
	return &PGStore{db: db}
}

func InitSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}

func (s *PGStore) Create(name, email, hashedPassword string) (models.User, error) {
	var id int64
	err := s.db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		name, email, hashedPassword,
	).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.User{}, ErrEmailExists
		}
		return models.User{}, fmt.Errorf("insert user: %w", err)
	}

	return models.User{
		ID:       fmt.Sprintf("%d", id),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}, nil
}

func (s *PGStore) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, name, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find by email: %w", err)
	}
	return user, nil
}

func (s *PGStore) FindByID(id string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, name, email, password FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return user, nil
}
