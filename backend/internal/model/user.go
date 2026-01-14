package model

type User struct {
	ID           int64  `db:"id"`
	FullName     string `db:"full_name"`
	Email        string `db:"email"`
	PasswordHash string `db:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
