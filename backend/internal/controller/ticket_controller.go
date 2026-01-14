package controller

import (
	"net/http"

	ticket_service "cinema.com/demo/internal/service/ticket"
	"github.com/gin-gonic/gin"
)

type GetTicketByUserIDRequest struct {
	UserID int `form:"user_id" binding:"required,gte=1"`
}

type TicketController struct {
	ticketService *ticket_service.TicketService
}

func NewTicketController(ticketService *ticket_service.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

func (c *TicketController) GetTicketByUserID(ctx *gin.Context) {
	var req GetTicketByUserIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tickets, err := c.ticketService.GetTicketByUserID(ctx.Request.Context(), req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tickets)
}
