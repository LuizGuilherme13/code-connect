package models

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type RegisterRequest struct {
	Name     string `json:"name"     example:"João"`
	Email    string `json:"email"    example:"joao@email.com"`
	Password string `json:"password" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email"    example:"joao@email.com"`
	Password string `json:"password" example:"123456"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"mensagem de erro"`
}

type UserResponse struct {
	ID    string `json:"id"    example:"1"`
	Name  string `json:"name"  example:"João"`
	Email string `json:"email" example:"joao@email.com"`
}
