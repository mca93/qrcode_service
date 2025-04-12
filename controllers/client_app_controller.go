package controllers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
	"github.com/mca93/qrcode_service/validators"

	"github.com/gin-gonic/gin"
)

// POST /v1/clientapps
func CreateClientApp(c *gin.Context) {
	var req models.ClientAppCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := validators.ValidateClientAppCreate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientApp := models.ClientApp{
		ID:           uuid.NewString(),
		Name:         req.Name,
		ContactEmail: req.ContactEmail,
		Status:       models.ClientAppStatusActive,
		CreatedAt:    time.Now(),
	}

	if err := config.DB.Create(&clientApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create client app"})
		return
	}

	c.JSON(http.StatusOK, clientApp)
}

// GET /v1/clientapps
func ListClientApps(c *gin.Context) {
	var clientApps []models.ClientApp
	if err := config.DB.Find(&clientApps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch client apps"})
		return
	}

	c.JSON(http.StatusOK, clientApps)
}

// GET /v1/clientapps/:id
func GetClientApp(c *gin.Context) {
	clientAppID := c.Param("id")
	var clientApp models.ClientApp

	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client app not found"})
		return
	}

	c.JSON(http.StatusOK, clientApp)
}

// PUT /v1/clientapps/:id
func UpdateClientApp(c *gin.Context) {
	clientAppID := c.Param("id")
	var req models.ClientAppUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := validators.ValidateClientAppUpdate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var clientApp models.ClientApp

	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client app not found"})
		return
	}

	if req.Name != "" {
		clientApp.Name = req.Name
	}
	if req.ContactEmail != "" {
		clientApp.ContactEmail = req.ContactEmail
	}
	if req.Status != "" {
		clientApp.Status = req.Status
	}

	if err := config.DB.Save(&clientApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update client app"})
		return
	}

	c.JSON(http.StatusOK, clientApp)
}
