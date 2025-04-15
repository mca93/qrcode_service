package controllers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
	"github.com/mca93/qrcode_service/validators"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// ListQRCodes retrieves all QR codes for the authenticated client app.
func ListQRCodes(c *gin.Context) {
	clientAppID := c.GetHeader("client_app_id")

	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientAppId header is required"})
		return
	}
	// Validate the client app ID
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client app not found"})
		return
	}
	// Check if the client app is active
	if clientApp.Status != models.ClientAppStatusActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Client app is not active"})
		return
	}
	var codes []models.QRCode
	if err := config.DB.Where("client_app_id = ?", clientAppID).Find(&codes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve QR codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_codes": codes})
}

// CreateQRCode handles the creation of a new QR code.
func CreateQRCode(c *gin.Context) {
	clientAppID, err := getClientAppID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.QRCodeCreateRequest
	req.ClientAppID = clientAppID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the template to validate the Data field
	var template models.Template
	if err := config.DB.First(&template, "id = ?", req.TemplateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Validate the Data field
	if err := validators.ValidateQRCodeData(req.Data, template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the QR code
	qrCode := models.QRCode{
		ID:            uuid.NewString(),
		Type:          req.Type,
		CreatedAt:     time.Now(),
		ExpiresAt:     req.ExpiresAt,
		Status:        "ACTIVE",
		ScanCount:     0,
		ImageURL:      "", // This should be generated later
		DeepLinkURL:   "", // This should be generated later
		ClientAppID:   req.ClientAppID,
		TemplateID:    req.TemplateID,
		ThirdPartyRef: req.ThirdPartyRef,
		Data:          req.Data,
	}

	// Save the QR code to the database
	if err := config.DB.Create(&qrCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create QR code"})
		return
	}

	c.JSON(http.StatusOK, qrCode)
}

// GetQRCode retrieves a specific QR code by its ID.
func GetQRCode(c *gin.Context) {
	clientAppID := c.GetHeader("client_app_id")
	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientAppId header is required"})
		return
	}

	id := c.Param("id")
	var qrCode models.QRCode
	if err := config.DB.First(&qrCode, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR Code not found"})
		return
	}

	// Validate ownership
	if err := validators.ValidateQRCodeOwnership(clientAppID, qrCode); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrCode)
}

// UpdateQRCode updates an existing QR code by its ID.
func UpdateQRCode(c *gin.Context) {
	clientAppID := c.GetHeader("client_app_id")
	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientAppId header is required"})
		return
	}

	id := c.Param("id")
	var qrCode models.QRCode
	if err := config.DB.First(&qrCode, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR Code not found"})
		return
	}

	// Validate ownership
	if err := validators.ValidateQRCodeOwnership(clientAppID, qrCode); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var req models.QRCodeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the update request
	if err := validators.ValidateQRCodeUpdate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply updates
	if req.Type != "" {
		qrCode.Type = req.Type
	}
	if req.ExpiresAt != nil {
		qrCode.ExpiresAt = req.ExpiresAt
	}
	if req.Status != "" {
		qrCode.Status = req.Status
	}
	if req.Data != nil {
		qrCode.Data = req.Data
	}

	config.DB.Save(&qrCode)
	c.JSON(http.StatusOK, qrCode)
}

// DeleteQRCode deletes a QR code by its ID.
func DeleteQRCode(c *gin.Context) {
	clientAppID := c.GetHeader("client_app_id")
	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientAppId header is required"})
		return
	}

	id := c.Param("id")
	var qrCode models.QRCode
	if err := config.DB.First(&qrCode, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR Code not found"})
		return
	}

	// Validate ownership
	if err := validators.ValidateQRCodeOwnership(clientAppID, qrCode); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	config.DB.Delete(&qrCode)
	c.JSON(http.StatusOK, gin.H{"message": "QR Code deleted successfully"})
}

// GetQRCodeImage generates and returns the QR code image.
func GetQRCodeImage(c *gin.Context) {
	id := c.Param("id")
	var qr models.QRCode

	if err := config.DB.First(&qr, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR Code not found"})
		return
	}

	png, err := qrcode.Encode(qr.DeepLinkURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

// ScanQRCode increments the scan count for a QR code.
func ScanQRCode(c *gin.Context) {
	id := c.Param("id")
	var qr models.QRCode

	if err := config.DB.First(&qr, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR Code not found"})
		return
	}

	qr.ScanCount++
	config.DB.Save(&qr)
	c.JSON(http.StatusOK, gin.H{"message": "Scan count updated", "scan_count": qr.ScanCount})
}
