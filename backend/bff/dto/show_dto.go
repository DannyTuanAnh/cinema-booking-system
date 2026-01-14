package dto

type GetShowQueryDTO struct {
	MovieID int `form:"movie_id" binding:"required,gte=1"`
}
