package controller

import (
	"net/http"

	"cinema.com/demo/bff/utils"
	seat_service "cinema.com/demo/internal/service/seat"
	"github.com/gin-gonic/gin"
)

type GetSeatByShowIDRequest struct {
	ShowID int `form:"show_id" binding:"required,gte=1"`
}

type SeatController struct {
	seatService *seat_service.SeatService
}

func NewSeatController(seatService *seat_service.SeatService) *SeatController {
	return &SeatController{seatService: seatService}
}

func (s *SeatController) GetSeatByShowID(ctx *gin.Context) {
	var req GetSeatByShowIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	seats, err := s.seatService.GetSeatByShowID(ctx.Request.Context(), req.ShowID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, seats)
}
