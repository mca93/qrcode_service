package controllers

import (
	"net/http"

	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"

	"github.com/gin-gonic/gin"
)

func ListApiKeys(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	var keys []models.ApiKey
	config.DB.Where("client_app_id = ?", clientAppID).Find(&keys)
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func CreateApiKey(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	var input models.ApiKey
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ClientAppID = clientAppID
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}
