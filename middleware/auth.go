package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
)

func QRCodeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.GetHeader("client_app_id")

		if appID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_app_id headers are required"})
			c.Abort()
			return
		}

		// Validate apiKeyID and appID in the database
		var clientApp models.ClientApp
		if err := config.DB.Where("id = ?", appID).First(&clientApp).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client_app_id"})
			c.Abort()
			return
		}

		// Check if the API key is active
		if clientApp.Status != models.ClientAppStatusActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Client app is not active!"})
			c.Abort()
			return
		}

		// If valid, proceed to the next handler
		c.Next()
	}
}
