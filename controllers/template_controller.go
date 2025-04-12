package controllers

import (
	"net/http"

	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"

	"github.com/gin-gonic/gin"
)

func ListTemplates(c *gin.Context) {
	var templates []models.Template
	config.DB.Find(&templates)
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func CreateTemplate(c *gin.Context) {
	var input models.Template
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}
