package middleware

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

func ApiKeyAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" || apiKey != os.Getenv("API_KEY") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()
    }
}
