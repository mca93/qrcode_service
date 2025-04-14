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

// ListApiKeys retrieves all API keys for a specific ClientAppID.
func ListApiKeys(c *gin.Context) {
	clientAppID := c.Param("id")
	status := c.DefaultQuery("status", string(models.ApiKeyStatusUnspecified))

	var keys []models.ApiKey
	query := config.DB.Where("client_app_id = ?", clientAppID)
	if status != string(models.ApiKeyStatusUnspecified) {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API keys"})
		return
	}

	var responses []models.ApiKeyResponse
	for _, key := range keys {
		responses = append(responses, toApiKeyResponse(key))
	}

	c.JSON(http.StatusOK, models.ApiKeyListResponse{
		ApiKeys:    responses,
		TotalCount: len(responses),
		Page:       1,
		PageSize:   len(responses),
		ItemsCount: len(responses),
		Items:      responses,
	})
}

// CreateApiKey creates a new API key for a specific ClientAppID.
func CreateApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	var req models.ApiKeyCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if ClientAppID exists in the database
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ClientAppID does not exist"})
		return
	}

	// Validate the request
	if _, err := validators.ValidateApiKeyCreate(req, clientAppID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := models.ApiKey{
		ID:          uuid.NewString(),
		Name:        req.Name,
		ClientAppID: clientAppID,
		KeyPrefix:   req.Name + "_" + uuid.NewString(),
		Status:      models.ApiKeyStatusActive,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(http.StatusOK, toApiKeyResponse(apiKey))
}

// GetApiKey retrieves a specific API key by its ID.
func GetApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Check if the requesting ClientAppID matches the API key's ClientAppID
	if key.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this API key"})
		return
	}

	c.JSON(http.StatusOK, toApiKeyResponse(key))
}

// UpdateApiKey updates an existing API key by its ID.
func UpdateApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var req models.ApiKeyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Check if the requesting ClientAppID matches the API key's ClientAppID
	if key.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this API key"})
		return
	}

	// Update fields
	if req.Name != "" {
		key.Name = req.Name
	}
	if req.Status != "" {
		key.Status = req.Status
	}
	key.UpdatedAt = time.Now()

	if err := config.DB.Save(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
		return
	}

	c.JSON(http.StatusOK, toApiKeyResponse(key))
}

// DeleteApiKey deletes an API key by its ID.
func DeleteApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Check if the requesting ClientAppID matches the API key's ClientAppID
	if key.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this API key"})
		return
	}

	if err := config.DB.Delete(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

// RegenerateApiKey regenerates the key prefix for an API key.
func RegenerateApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Check if the requesting ClientAppID matches the API key's ClientAppID
	if key.ClientAppID != clientAppID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to regenerate this API key"})
		return
	}

	// Generate a new key prefix
	key.KeyPrefix = key.Name + "_" + uuid.NewString()
	key.UpdatedAt = time.Now()

	if err := config.DB.Save(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to regenerate API key"})
		return
	}

	c.JSON(http.StatusOK, toApiKeyResponse(key))
}

// Helper to convert to ApiKeyResponse
func toApiKeyResponse(key models.ApiKey) models.ApiKeyResponse {
	return models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		Name:        key.Name,
		KeyPrefix:   key.KeyPrefix,
		Status:      key.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	}
}
