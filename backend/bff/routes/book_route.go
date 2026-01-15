package routes

import (
	"database/sql"
	"os"

	"cinema.com/demo/bff/clients/book"
	"cinema.com/demo/bff/controllers/book"
	"cinema.com/demo/bff/middleware"
	"github.com/gin-gonic/gin"
)

func InitBookRoutes(r *gin.RouterGroup, db *sql.DB) {
	bookGroup := r.Group("/book")
	RegisterBookRoutes(bookGroup, db)
}

func RegisterBookRoutes(r *gin.RouterGroup, db *sql.DB) {
	addr := os.Getenv("ADDR_SERVER")
	path := "https://" + addr + "/api"

	bookClient := book_clients.NewBookHTTPClient(path)

	bookController := book.NewBookController(bookClient)

	RegisterBookRoute(r, bookController, db)
}

func RegisterBookRoute(r *gin.RouterGroup, b *book.BookController, db *sql.DB) {
	r.POST("", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), b.Book)
}
