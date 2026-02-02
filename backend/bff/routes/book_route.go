package routes

import (
	"os"

	"cinema.com/demo/bff/clients/book"
	"cinema.com/demo/bff/controllers/book"
	"github.com/gin-gonic/gin"
)

func InitBookRoutes(r *gin.RouterGroup) {
	bookGroup := r.Group("/book")
	RegisterBookRoutes(bookGroup)
}

func RegisterBookRoutes(r *gin.RouterGroup) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	bookClient := book_clients.NewBookHTTPClient(path)

	bookController := book.NewBookController(bookClient)

	RegisterBookRoute(r, bookController)
}

func RegisterBookRoute(r *gin.RouterGroup, b *book.BookController) {
	r.POST("", b.Book)
}
