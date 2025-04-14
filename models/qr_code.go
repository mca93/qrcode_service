package models

import (
	"time"
)

// QRCodeType defines the allowed types for QR codes.
type QRCodeType string

const (
	QRCodeTypeStable  QRCodeType = "STABLE"
	QRCodeTypeDynamic QRCodeType = "DYNAMIC"
)

// QRCode represents the QR code entity.
type QRCode struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	Type         QRCodeType `gorm:"not null" json:"type"` // Restricted to STABLE or DYNAMIC
	CreatedAt    time.Time  `json:"createdAt"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	Status       string     `json:"status"`
	ScanCount    int64      `json:"scanCount"`
	ImageURL     string     `json:"imageUrl"`
	DeepLinkURL  string     `json:"deepLinkUrl"`
	ClientAppID  string     `gorm:"not null" json:"clientAppId"`                   // Foreign key to ClientApp
	TemplateID   string     `gorm:"not null" json:"templateId"`                    // Foreign key to Template
	ThirdPartRef string     `json:"thirdPartRef"`                                  // Reference to third-party systems
	Data         JSONMap    `gorm:"type:jsonb" json:"data"`                        // Custom key-value data
	ClientApp    ClientApp  `gorm:"foreignKey:ClientAppID;references:ID" json:"-"` // Association with ClientApp
	Template     Template   `gorm:"foreignKey:TemplateID;references:ID" json:"-"`  // Association with Template
}

// QRCodeCreateRequest represents the request structure for creating a QR code.
type QRCodeCreateRequest struct {
	Type         QRCodeType             `json:"type" binding:"required,oneof=STABLE DYNAMIC"` // Restricted to STABLE or DYNAMIC
	ExpiresAt    *time.Time             `json:"expiresAt,omitempty"`
	TemplateID   string                 `json:"templateId" binding:"required"`
	ClientAppID  string                 `json:"clientAppId" binding:"required"`
	ThirdPartRef string                 `json:"thirdPartRef"`
	Data         map[string]interface{} `json:"data"` // Custom key-value data
}

// QRCodeUpdateRequest represents the request structure for updating a QR code.
type QRCodeUpdateRequest struct {
	Type         QRCodeType             `json:"type,omitempty" binding:"omitempty,oneof=STABLE DYNAMIC"` // Restricted to STABLE or DYNAMIC
	ExpiresAt    *time.Time             `json:"expiresAt,omitempty"`
	Status       string                 `json:"status,omitempty"`
	ThirdPartRef string                 `json:"thirdPartRef,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"` // Custom key-value data
}

// QRCodeResponse represents the response structure for a QR code.
type QRCodeResponse struct {
	ID           string                 `json:"id"`
	Type         QRCodeType             `json:"type"`
	CreatedAt    time.Time              `json:"createdAt"`
	ExpiresAt    *time.Time             `json:"expiresAt,omitempty"`
	Status       string                 `json:"status"`
	ScanCount    int64                  `json:"scanCount"`
	ImageURL     string                 `json:"imageUrl"`
	DeepLinkURL  string                 `json:"deepLinkUrl"`
	ClientAppID  string                 `json:"clientAppId"`
	TemplateID   string                 `json:"templateId"`
	ThirdPartRef string                 `json:"thirdPartRef"`
	Data         map[string]interface{} `json:"data"` // Custom key-value data
}
