package controller

import (
	"net/http"

	"cinema.com/demo/internal/domain"
	"cinema.com/demo/internal/model"
	book_service "cinema.com/demo/internal/service/book"
	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookService *book_service.BookService
}

func NewBookController(bookService *book_service.BookService) *BookController {
	return &BookController{bookService: bookService}
}

func (b *BookController) Book(ctx *gin.Context) {
	var req model.BookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := b.bookService.BookSeats(ctx.Request.Context(), req.UserID, req.Seats)
	if err != nil {

		if err == domain.ErrSeatAlreadyBooked {
			ctx.JSON(http.StatusConflict, gin.H{"error": "One or more seats are already booked"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book seats"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Seats booked successfully"})
}
