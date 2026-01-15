package model

type BookRequest struct {
	UserID int   `json:"user_id" binding:"required,gte=1"`
	Seats  []int `json:"seats" binding:"required,min=1,dive,gte=1"`
}
