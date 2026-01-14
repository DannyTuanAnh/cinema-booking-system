package ticket

import (
	"net/http"

	"cinema.com/demo/bff/clients/ticket"
	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketClient ticket_clients.TicketClient
}

func NewTicketController(ticketClient ticket_clients.TicketClient) *TicketController {
	return &TicketController{
		ticketClient: ticketClient,
	}
}

func (t *TicketController) GetTicketByUserID(ctx *gin.Context) {

	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	// convert userId type any to int
	userIdInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}

	resp, err := t.ticketClient.GetTicketByUserID(userIdInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": resp})

}
