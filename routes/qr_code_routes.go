package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/controllers"
	"github.com/mca93/qrcode_service/middleware"
)

func RegisterQRCodeRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		// QR code routes with middleware to validate API key and ClientAppID
		v1.GET("/qrcodes", middleware.QRCodeAuthMiddleware(), controllers.ListQRCodes)
		v1.POST("/qrcodes", middleware.QRCodeAuthMiddleware(), controllers.CreateQRCode)
		v1.GET("/qrcodes/:id", middleware.QRCodeAuthMiddleware(), controllers.GetQRCode)
		// Preview QR code image
		v1.GET("/qrcodes/:id/preview", controllers.GetQRCodeImage)
		// Download QR code image
		// v1.GET("/qrcodes/:id/download", middleware.ApiKeyAuthMiddleware(), controllers.DownloadQRCode)
		v1.PUT("/qrcodes/:id", middleware.QRCodeAuthMiddleware(), controllers.UpdateQRCode)
		v1.DELETE("/qrcodes/:id", middleware.QRCodeAuthMiddleware(), controllers.DeleteQRCode)
	}
}
