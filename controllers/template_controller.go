package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
	"github.com/mca93/qrcode_service/validators"
	"gorm.io/gorm"
)

type TemplateController struct {
	// We'll keep using config.DB directly as per your existing setup
}

func NewTemplateController() *TemplateController {
	return &TemplateController{}
}

// CreateTemplate handles the creation of a new template.
func (tc *TemplateController) CreateTemplate(c *gin.Context) {
	clientAppID, err := getValidClientAppID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // Limit to 10 MB
		respondWithError(c, http.StatusBadRequest, "Failed to parse multipart form data")
		return
	}

	// Extract form fields
	name := c.PostForm("name")
	description := c.PostForm("description")
	shape := c.PostForm("shape")
	foregroundColor := c.PostForm("foregroundColor")
	backgroundColor := c.PostForm("backgroundColor")
	size := c.PostForm("size")
	errorCorrection := c.PostForm("errorCorrection")
	definitionJSON := c.PostForm("definition")

	// Validate required fields
	if name == "" || clientAppID == "" {
		respondWithError(c, http.StatusBadRequest, "Name and ClientAppID are required")
		return
	}

	// Parse and validate the definition field (JSON)
	var definition models.Definition
	if err := json.Unmarshal([]byte(definitionJSON), &definition); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid definition format")
		return
	}
	if err := definition.Validate(); err != nil {
		respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Invalid definition: %v", err))
		return
	}

	// Handle logo file upload
	logoPath, err := handleLogoUpload(c)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to process logo upload")
		return
	}

	// Convert size and errorCorrection to appropriate types
	sizeInt, err := strconv.Atoi(size)
	if err != nil || sizeInt <= 0 {
		respondWithError(c, http.StatusBadRequest, "Invalid size value")
		return
	}

	// Map errorCorrection to enum
	errorCorrectionEnum := models.QRCodeErrorCorrection(errorCorrection)
	if errorCorrectionEnum == "" {
		respondWithError(c, http.StatusBadRequest, "Invalid errorCorrection value")
		return
	}

	// Create the template
	template := models.Template{
		ID:              uuid.NewString(),
		Name:            name,
		Description:     description,
		ClientAppID:     clientAppID,
		Definition:      definition,
		Shape:           shape,
		ForegroundColor: foregroundColor,
		BackgroundColor: backgroundColor,
		Size:            sizeInt,
		LogoURL:         logoPath,
		ErrorCorrection: errorCorrectionEnum,
		Active:          true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save the template to the database
	if err := config.DB.Create(&template).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create template")
		return
	}

	respondWithSuccess(c, http.StatusCreated, template)
}

// ListTemplates retrieves all templates for a specific ClientAppID.
func (tc *TemplateController) ListTemplates(c *gin.Context) {
	clientAppID, err := getValidClientAppID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var templates []models.Template
	if err := config.DB.Where("client_app_id = ?", clientAppID).Find(&templates).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve templates")
		return
	}

	respondWithSuccess(c, http.StatusOK, templates)
}

// GetTemplate retrieves a specific template by its ID.
func (tc *TemplateController) GetTemplate(c *gin.Context) {
	clientAppID, err := getValidClientAppID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")
	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondWithError(c, http.StatusNotFound, "Template not found")
		} else {
			respondWithError(c, http.StatusInternalServerError, "Failed to retrieve template")
		}
		return
	}

	if template.ClientAppID != clientAppID {
		respondWithError(c, http.StatusForbidden, "You do not have permission to access this template")
		return
	}

	respondWithSuccess(c, http.StatusOK, template)
}

// UpdateTemplate updates an existing template by its ID.
func (tc *TemplateController) UpdateTemplate(c *gin.Context) {
	clientAppID, err := getValidClientAppID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")
	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondWithError(c, http.StatusNotFound, "Template not found")
		} else {
			respondWithError(c, http.StatusInternalServerError, "Failed to retrieve template")
		}
		return
	}

	if template.ClientAppID != clientAppID {
		respondWithError(c, http.StatusForbidden, "You do not have permission to update this template")
		return
	}

	var req models.TemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validators.ValidateTemplateUpdate(req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	logoPath, err := handleLogoUpload(c)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to process logo upload")
		return
	}

	// Update template fields
	template.Name = req.Name
	template.Description = req.Description
	template.Definition = req.Definition
	template.Shape = req.Shape
	template.ForegroundColor = req.ForegroundColor
	template.BackgroundColor = req.BackgroundColor
	template.Size = req.Size
	template.ErrorCorrection = req.ErrorCorrection
	template.UpdatedAt = time.Now()

	if logoPath != "" {
		// Delete old logo if exists
		if template.LogoURL != "" {
			if err := deleteFile(template.LogoURL); err != nil {
				// Log the error but continue with the update
				fmt.Printf("Failed to delete old logo: %v\n", err)
			}
		}
		template.LogoURL = logoPath
	}

	if err := config.DB.Save(&template).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update template")
		return
	}

	respondWithSuccess(c, http.StatusOK, template)
}

// DeactivateTemplate deactivates a template by its ID.
func (tc *TemplateController) DeactivateTemplate(c *gin.Context) {
	clientAppID, err := getValidClientAppID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")
	var template models.Template
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondWithError(c, http.StatusNotFound, "Template not found")
		} else {
			respondWithError(c, http.StatusInternalServerError, "Failed to retrieve template")
		}
		return
	}

	if template.ClientAppID != clientAppID {
		respondWithError(c, http.StatusForbidden, "You do not have permission to deactivate this template")
		return
	}

	template.Active = false
	template.UpdatedAt = time.Now()

	if err := config.DB.Save(&template).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to deactivate template")
		return
	}

	respondWithSuccess(c, http.StatusOK, template)
}

// Helper functions (could be moved to a separate package if needed)

func getValidClientAppID(c *gin.Context) (string, error) {
	clientAppID := c.GetHeader("client_app_id")
	if clientAppID == "" {
		return "", errors.New("ClientAppID header is required")
	}

	// Validate if ClientAppID exists in the database
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		return "", errors.New("ClientAppID does not exist")
	}

	return clientAppID, nil
}

func handleLogoUpload(c *gin.Context) (string, error) {
	file, err := c.FormFile("logo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil // No file uploaded is not an error
		}
		return "", err
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%s", uuid.NewString(), file.Filename)
	uploadDir := getEnv("UPLOAD_DIR", "./uploads/logos")
	logoPath := filepath.Join(uploadDir, filename)

	// Ensure the directory exists
	if err := ensureDirectoryExists(uploadDir); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save the uploaded file
	if err := c.SaveUploadedFile(file, logoPath); err != nil {
		return "", fmt.Errorf("failed to save logo file: %w", err)
	}

	return logoPath, nil
}

func deleteFile(path string) error {
	if path == "" {
		return nil
	}
	return os.Remove(path)
}

func ensureDirectoryExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func respondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
