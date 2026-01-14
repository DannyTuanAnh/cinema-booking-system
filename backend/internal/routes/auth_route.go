package routes

import (
	"cinema.com/demo/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup, ctrl *controller.AuthController) {
	r.POST("/auth/login", ctrl.Login)
	r.POST("/auth/register", ctrl.Register)
}
