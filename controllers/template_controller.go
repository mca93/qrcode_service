package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
	"github.com/mca93/qrcode_service/validators"
)

func CreateTemplate(c *gin.Context) {
	var req models.TemplateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validators.ValidateTemplateCreate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := models.Template{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		ClientAppID: req.ClientAppID,
		Style:       req.Style,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func ListTemplates(c *gin.Context) {
	clientAppID := c.Query("clientAppId")
	var templates []models.Template

	query := config.DB
	if clientAppID != "" {
		query = query.Where("client_app_id = ?", clientAppID)
	}
	query.Find(&templates)

	c.JSON(http.StatusOK, templates)
}

func GetTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

func UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var req models.TemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validators.ValidateTemplateUpdate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Atualizações
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Style != nil {
		template.Style = *req.Style
	}
	template.UpdatedAt = time.Now()

	config.DB.Save(&template)
	c.JSON(http.StatusOK, template)
}

func DeactivateTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	template.Active = false
	template.UpdatedAt = time.Now()

	config.DB.Save(&template)
	c.JSON(http.StatusOK, template)
}
