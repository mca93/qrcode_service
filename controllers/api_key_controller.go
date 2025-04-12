package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
)

// GET /v1/clientapps/:clientAppId/apikeys
func ListApiKeys(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
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

// POST /v1/clientapps/:clientAppId/apikeys
func CreateApiKey(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	var req models.ApiKeyCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := models.ApiKey{
		ID:          uuid.NewString(),
		ClientAppID: clientAppID,
		KeyPrefix:   req.KeyPrefix,
		Status:      req.Status,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API Key"})
		return
	}

	response := models.ApiKeyResponse{
		ID:          apiKey.ID,
		ClientAppID: apiKey.ClientAppID,
		KeyPrefix:   apiKey.KeyPrefix,
		Status:      apiKey.Status,
		CreatedAt:   apiKey.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GET /v1/clientapps/:clientAppId/apikeys/:keyId
func GetApiKey(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	response := models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		KeyPrefix:   key.KeyPrefix,
		Status:      key.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	}

	c.JSON(http.StatusOK, response)
}

// PUT /v1/clientapps/:clientAppId/apikeys/:keyId
func UpdateApiKey(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	keyID := c.Param("keyId")

	var req models.ApiKeyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	config.DB.Model(&key).Updates(models.ApiKey{
		KeyPrefix:   req.KeyPrefix,
		ClientAppID: req.ClientAppID,
		Status:      req.Status,
	})

	c.JSON(http.StatusOK, models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		KeyPrefix:   req.KeyPrefix,
		Status:      req.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	})
}

// DELETE /v1/clientapps/:clientAppId/apikeys/:keyId
func DeleteApiKey(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	keyID := c.Param("keyId")

	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).Delete(&models.ApiKey{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API Key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API Key deleted successfully"})
}

// POST /v1/clientapps/:clientAppId/apikeys/:keyId/regenerate
func RegenerateApiKey(c *gin.Context) {
	clientAppID := c.Param("clientAppId")
	keyID := c.Param("keyId")

	var key models.ApiKey
	if err := config.DB.Where("client_app_id = ? AND id = ?", clientAppID, keyID).First(&key).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	newPrefix := generateNewKey()
	key.KeyPrefix = newPrefix

	if err := config.DB.Save(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to regenerate API Key"})
		return
	}

	c.JSON(http.StatusOK, models.ApiKeyResponse{
		ID:          key.ID,
		ClientAppID: key.ClientAppID,
		KeyPrefix:   key.KeyPrefix,
		Status:      key.Status,
		CreatedAt:   key.CreatedAt,
		RevokedAt:   key.RevokedAt,
	})
}

func generateNewKey() string {
	// Gerador simples — substitua por algo mais seguro se necessário
	return "key_" + uuid.NewString()[:8]
}
