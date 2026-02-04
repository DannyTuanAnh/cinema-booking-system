package seat

import (
	"net/http"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/seat"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/dto"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/utils"
	"github.com/gin-gonic/gin"
)

type SeatController struct {
	seatClient seat_clients.SeatClient
}

func NewSeatController(seatClient seat_clients.SeatClient) *SeatController {
	return &SeatController{
		seatClient: seatClient,
	}
}

func (s *SeatController) GetSeatByShowID(ctx *gin.Context) {
	var param dto.GetSeatQueryDTO
	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	resp, err := s.seatClient.GetSeatByShowID(param.ShowID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": resp})
}
