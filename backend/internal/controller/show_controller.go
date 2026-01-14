package controller

import (
	"net/http"

	show_service "cinema.com/demo/internal/service/show"
	"github.com/gin-gonic/gin"
)

type GetShowByMovieIDRequest struct {
	MovieID int `form:"movie_id" binding:"required,gte=1"`
}

type ShowController struct {
	showService *show_service.ShowService
}

func NewShowController(showService *show_service.ShowService) *ShowController {
	return &ShowController{showService: showService}
}

func (c *ShowController) GetShowByMovieID(ctx *gin.Context) {
	var req GetShowByMovieIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shows, err := c.showService.GetShowByMovieID(ctx.Request.Context(), req.MovieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, shows)
}
