package dto

type GetSeatQueryDTO struct {
	ShowID int `form:"show_id" binding:"required,gte=1"`
}
