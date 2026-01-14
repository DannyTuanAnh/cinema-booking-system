package dto

type BookSeatDTO struct {
	Seats []int `json:"seats" binding:"required,min=1,dive,gte=1"`
}
