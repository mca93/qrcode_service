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

// CreateTemplate handles the creation of a new template.
func CreateTemplate(c *gin.Context) {
	var req models.TemplateCreateRequest
	clientAppID := c.GetHeader("client_app_id")

	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ClientAppID is required"})
		return
	}

	// Validate if ClientAppID exists in the database
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ClientAppID does not exist"})
		return
	}

	req.ClientAppID = clientAppID

	// Bind the JSON request to the TemplateCreateRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request
	if err := validators.ValidateTemplateCreate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the template
	template := models.Template{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		ClientAppID: req.ClientAppID,
		Metadata:    req.Metadata,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save the template to the database
	if err := config.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	// Return the created template
	c.JSON(http.StatusOK, template)
}

// ListTemplates retrieves all templates for a specific ClientAppID.
func ListTemplates(c *gin.Context) {
	clientAppID := c.GetHeader("client_app_id")
	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ClientAppID is required"})
		return
	}

	// Validate if ClientAppID exists in the database
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ClientAppID is invalid"})
		return
	}

	var templates []models.Template
	if err := config.DB.Where("client_app_id = ?", clientAppID).Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetTemplate retrieves a specific template by its ID.
func GetTemplate(c *gin.Context) {
	id := c.Param("id")
	clientAppID := c.GetHeader("client_app_id")

	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ClientAppID header is required"})
		return
	}

	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Check if the requesting ClientAppID matches the template's ClientAppID
	if template.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate updates an existing template by its ID.
func UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	clientAppID := c.GetHeader("client_app_id")

	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_app_id header is required"})
		return
	}

	var req models.TemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Check if the requesting ClientAppID matches the template's ClientAppID
	if template.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this template"})
		return
	}

	// Validate the update request
	if req.Metadata != nil {
		if err := validators.ValidateTemplateUpdate(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Apply updates
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Metadata != nil {
		template.Metadata = *req.Metadata
	}
	template.UpdatedAt = time.Now()

	// Save the updated template
	if err := config.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeactivateTemplate deactivates a template by its ID.
func DeactivateTemplate(c *gin.Context) {
	id := c.Param("id")
	clientAppID := c.GetHeader("client_app_id")

	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_app_id header is required"})
		return
	}

	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Check if the requesting ClientAppID matches the template's ClientAppID
	if template.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to deactivate this template"})
		return
	}

	template.Active = false
	template.UpdatedAt = time.Now()

	// Save the deactivated template
	if err := config.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate template"})
		return
	}

	c.JSON(http.StatusOK, template)
}
