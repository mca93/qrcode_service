package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/controllers"
)

func RegisterClientAppRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.POST("/clientapps", controllers.CreateClientApp)
		v1.GET("/clientapps", controllers.ListClientApps)
		v1.GET("/clientapps/:id", controllers.GetClientApp)    // Obtenha um ClientApp específico pelo ID
		v1.PUT("/clientapps/:id", controllers.UpdateClientApp) // Atualize um ClientApp específico pelo ID
		// outras rotas podem ser adicionadas aqui
	}
}
