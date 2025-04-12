package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/controllers"
)

func RegisterClientAppRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.POST("/clientapps", controllers.CreateClientApp)
		// outras rotas podem ser adicionadas aqui
	}
}
