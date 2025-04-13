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

// GET /v1/clientapps/:id/apikeys
func ListApiKeys(c *gin.Context) {
	clientAppID := c.Param("id")
	status := c.DefaultQuery("status", string(models.ApiKeyStatusUnspecified))

	var keys []models.ApiKey
	query := config.DB.Where("client_app_id = ?", clientAppID)
	if status != string(models.ApiKeyStatusUnspecified) {
		query = query.Where("status = ?", status)
	}
	query.Find(&keys)

	var responses []models.ApiKeyResponse
	for _, key := range keys {
		responses = append(responses, models.ApiKeyResponse{
			ID:          key.ID,
			ClientAppID: key.ClientAppID,
			Name:        key.Name,
			KeyPrefix:   key.KeyPrefix,
			Status:      key.Status,
			CreatedAt:   key.CreatedAt,
			RevokedAt:   key.RevokedAt,
		})
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

// POST /v1/clientapps/:id/apikeys
func CreateApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	var req models.ApiKeyCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔐 Verifica se o clientApp existe
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clientAppId: client app does not exist"})
		return
	}

	// ✅ Validação e prefixo
	keyPrefix, err := validators.ValidateApiKeyCreate(req, clientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := models.ApiKey{
		ID:          uuid.NewString(),
		Name:        req.Name,
		ClientAppID: clientAppID,
		KeyPrefix:   keyPrefix,
		Status:      models.ApiKeyStatusActive,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API Key"})
		return
	}

	c.JSON(http.StatusOK, toApiKeyResponse(apiKey))
}

// GET /v1/clientapps/:id/apikeys/:keyId
func GetApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	c.JSON(http.StatusOK, models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		Name:        key.Name,
		KeyPrefix:   key.KeyPrefix,
		Status:      key.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	})
}

// PUT /v1/clientapps/:id/apikeys/:keyId
func UpdateApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var req models.ApiKeyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the update request
	if err := validators.ValidateApiKeyUpdate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	config.DB.Model(&key).Updates(models.ApiKey{
		Name:        req.Name,
		Status:      req.Status,
		ClientAppID: clientAppID, // garantido pela URL
	})

	c.JSON(http.StatusOK, models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		Name:        key.Name,
		KeyPrefix:   key.KeyPrefix,
		Status:      key.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	})
}

// DELETE /v1/clientapps/:id/apikeys/:keyId
func DeleteApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).Delete(&models.ApiKey{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API Key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API Key deleted successfully"})
}

// POST /v1/clientapps/:id/apikeys/:keyId/regenerate
func RegenerateApiKey(c *gin.Context) {
	clientAppID := c.Param("id")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	newPrefix := key.Name + "_" + clientAppID + "_" + time.Now().UTC().Format("20060102150405")
	key.KeyPrefix = newPrefix

	if err := config.DB.Save(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to regenerate API Key"})
		return
	}

	c.JSON(http.StatusOK, models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		Name:        key.Name,
		KeyPrefix:   key.KeyPrefix,
		Status:      key.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	})
}

// Helper para converter para ApiKeyResponse
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
