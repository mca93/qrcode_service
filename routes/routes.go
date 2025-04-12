package routes

import (
	"github.com/mca93/qrcode_service/controllers"
	"github.com/mca93/qrcode_service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1", middleware.ApiKeyAuthMiddleware())
	{
		v1.GET("/clientapps", controllers.ListClientApps)

		v1.GET("/clientapps/:clientAppId/apikeys", controllers.ListApiKeys)
		v1.POST("/clientapps/:clientAppId/apikeys", controllers.CreateApiKey)

		v1.GET("/qrcodes", controllers.ListQRCodes)
		v1.POST("/qrcodes", controllers.CreateQRCode)

		v1.GET("/templates", controllers.ListTemplates)
		v1.POST("/templates", controllers.CreateTemplate)
	}
	v1.POST("/clientapps", controllers.CreateClientApp)
	v1.GET("/qrcodes/:id/image", controllers.GetQRCodeImage)

}
