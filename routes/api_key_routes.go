package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/controllers"
)

func RegisterApiKeyRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		// Rotas aninhadas por ClientApp
		clientApp := v1.Group("/clientapps/:id/apikeys")
		{
			clientApp.GET("", controllers.ListApiKeys)
			clientApp.POST("", controllers.CreateApiKey)
			clientApp.GET("/:keyId", controllers.GetApiKey)
			clientApp.PUT("/:keyId", controllers.UpdateApiKey)
			clientApp.DELETE("/:keyId", controllers.DeleteApiKey)
			clientApp.POST("/:keyId/regenerate", controllers.RegenerateApiKey)
		}
	}
}
