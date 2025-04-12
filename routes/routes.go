package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// v1 := router.Group("/v1")
	// {
	// v1.GET("/clientapps", controllers.ListClientApps)
	// v1.POST("/clientapps", controllers.CreateClientApp)

	// Regista as rotas do ClientApp
	RegisterClientAppRoutes(router)
	// Regista as rotas do ApiKey
	RegisterApiKeyRoutes(router)

	// v1.GET("/qrcodes/:id/image", controllers.GetQRCodeImage)
	// v1.GET("/clientapps/:clientAppId/apikeys", controllers.ListApiKeys)
	// v1.POST("/clientapps/:clientAppId/apikeys", controllers.CreateApiKey)

	// v1.GET("/qrcodes", middleware.ApiKeyAuthMiddleware(), controllers.ListQRCodes)
	// v1.POST("/qrcodes", middleware.ApiKeyAuthMiddleware(), controllers.CreateQRCode)

	// v1.GET("/templates", middleware.ApiKeyAuthMiddleware(), controllers.ListTemplates)
	// v1.POST("/templates", middleware.ApiKeyAuthMiddleware(), controllers.CreateTemplate)
	//}

}
