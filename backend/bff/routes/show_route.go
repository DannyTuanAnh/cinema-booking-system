package routes

import (
	"database/sql"
	"os"

	"cinema.com/demo/bff/clients/show"
	"cinema.com/demo/bff/controllers/show"
	"cinema.com/demo/bff/middleware"
	"github.com/gin-gonic/gin"
)

func InitShowRoutes(r *gin.RouterGroup, db *sql.DB) {
	showGroup := r.Group("/shows")
	RegisterShowRoutes(showGroup, db)
}

func RegisterShowRoutes(r *gin.RouterGroup, db *sql.DB) {
	addr := os.Getenv("ADDR_SERVER")
	path := "http://" + addr + "/api"

	showClient := show_clients.NewShowHTTPClient(path)

	showController := show.NewShowController(showClient)

	RegisterGetShowRoute(r, showController, db)
}

func RegisterGetShowRoute(r *gin.RouterGroup, s *show.ShowController, db *sql.DB) {
	r.GET("", middleware.ApiKeyMiddleware(db), middleware.RateLimit(), s.GetShowByMovieID)
}
