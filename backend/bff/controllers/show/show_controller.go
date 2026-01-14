package show

import (
	"net/http"

	"cinema.com/demo/bff/clients/show"
	"cinema.com/demo/bff/dto"
	"cinema.com/demo/bff/utils"
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
