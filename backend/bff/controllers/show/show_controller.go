package show

import (
	"net/http"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/show"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/dto"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/utils"
	"github.com/gin-gonic/gin"
)

type ShowController struct {
	showClient show_clients.ShowClient
}

func NewShowController(showClient show_clients.ShowClient) *ShowController {
	return &ShowController{
		showClient: showClient,
	}
}

func (sc *ShowController) GetShowByMovieID(ctx *gin.Context) {
	var param dto.GetShowQueryDTO
	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	resp, err := sc.showClient.GetShowByMovieID(param.MovieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": resp})
}
