package book

import (
	"net/http"

	"github.com/DannyTuanAnh/cinema-booking-system/bff/clients/book"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/dto"
	"github.com/DannyTuanAnh/cinema-booking-system/bff/utils"
	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookClient book_clients.BookClient
}

func NewBookController(bookClient book_clients.BookClient) *BookController {
	return &BookController{
		bookClient: bookClient,
	}
}

func (b *BookController) Book(ctx *gin.Context) {
	var param dto.BookSeatDTO

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	// convert userId type any to int64
	userIdInt, ok := userId.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}

	err := b.bookClient.BookSeats(book_clients.BookRequest{
		UserID: userIdInt,
		Seats:  param.Seats,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": "booking successful"})

}
