package controllers

import (
	"net/http"

	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"

	"github.com/gin-gonic/gin"
)

func ListClientApps(c *gin.Context) {
	var apps []models.ClientApp
	config.DB.Find(&apps)
	c.JSON(http.StatusOK, gin.H{"apps": apps})
}

func CreateClientApp(c *gin.Context) {
	var input models.ClientApp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}
