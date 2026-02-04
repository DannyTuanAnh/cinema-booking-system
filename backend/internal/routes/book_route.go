package routes

import (
	"github.com/DannyTuanAnh/cinema-booking-system/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitBookRoutes(r *gin.RouterGroup, ctrl *controller.BookController) {
	r.POST("/book", ctrl.Book)
}
