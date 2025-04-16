package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

// QRCodeType defines the allowed types for QR codes.
type QRCodeType string

const (
	QRCodeTypeStable  QRCodeType = "STABLE"
	QRCodeTypeDynamic QRCodeType = "DYNAMIC"
)

// QRCode represents the QR code entity.
type QRCode struct {
	ID            string     `gorm:"primaryKey" json:"id"`
	Type          QRCodeType `gorm:"not null" json:"type"` // Restricted to STABLE or DYNAMIC
	CreatedAt     time.Time  `json:"createdAt"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	Status        string     `json:"status"`
	ScanCount     int64      `json:"scanCount"`
	ImageURL      string     `json:"imageUrl"`
	DeepLinkURL   string     `json:"deepLinkUrl"`                                   // Auto-generated deep link
	ClientAppID   string     `gorm:"not null" json:"clientAppId"`                   // Foreign key to ClientApp
	TemplateID    string     `gorm:"not null" json:"templateId"`                    // Foreign key to Template
	ThirdPartyRef string     `json:"thirdPartRef"`                                  // Reference to third-party systems
	Data          JSONMap    `gorm:"type:jsonb" json:"data"`                        // Custom key-value data
	ClientApp     ClientApp  `gorm:"foreignKey:ClientAppID;references:ID" json:"-"` // Association with ClientApp
	Template      Template   `gorm:"foreignKey:TemplateID;references:ID" json:"-"`  // Association with Template
}

// BeforeCreate is a GORM hook that runs before a new QRCode is inserted into the database.
func (q *QRCode) BeforeCreate(tx *gorm.DB) (err error) {
	// Auto-generate the DeepLinkURL if it is not already set
	if q.DeepLinkURL == "" {
		protocol := getEnv("DEEPLINK_PROTOCOL", "https")  // Default to "https" if not set
		host := getEnv("DEEPLINK_HOST", "yourdomain.com") // Default to "yourdomain.com" if not set
		q.DeepLinkURL = fmt.Sprintf("%s://%s/qrcodes/%s", protocol, host, q.ID)
	}
	return nil
}

// getEnv retrieves the value of an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// QRCodeCreateRequest represents the request structure for creating a QR code.
type QRCodeCreateRequest struct {
	Type          QRCodeType             `json:"type" binding:"required,oneof=STABLE DYNAMIC"` // Restricted to STABLE or DYNAMIC
	ExpiresAt     *time.Time             `json:"expiresAt,omitempty"`
	TemplateID    string                 `json:"templateId" binding:"required"`
	ClientAppID   string                 `json:"clientAppId" binding:"required"`
	ThirdPartyRef string                 `json:"third_party_ref"`
	Data          map[string]interface{} `json:"data"` // Custom key-value data
}

// QRCodeUpdateRequest represents the request structure for updating a QR code.
type QRCodeUpdateRequest struct {
	Type          QRCodeType             `json:"type,omitempty" binding:"omitempty,oneof=STABLE DYNAMIC"` // Restricted to STABLE or DYNAMIC
	ExpiresAt     *time.Time             `json:"expiresAt,omitempty"`
	Status        string                 `json:"status,omitempty"`
	ThirdPartyRef string                 `json:"third_party_ref,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"` // Custom key-value data
}

// QRCodeResponse represents the response structure for a QR code.
type QRCodeResponse struct {
	ID            string                 `json:"id"`
	Type          QRCodeType             `json:"type"`
	CreatedAt     time.Time              `json:"createdAt"`
	ExpiresAt     *time.Time             `json:"expiresAt,omitempty"`
	Status        string                 `json:"status"`
	ScanCount     int64                  `json:"scanCount"`
	ImageURL      string                 `json:"imageUrl"`
	DeepLinkURL   string                 `json:"deepLinkUrl"`
	ClientAppID   string                 `json:"clientAppId"`
	TemplateID    string                 `json:"templateId"`
	ThirdPartyRef string                 `json:"third_party_ref"`
	Data          map[string]interface{} `json:"data"` // Custom key-value data
}

// JSONMap is a custom type to handle JSON fields in the database.
type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]uint8)
	if !ok {
		return errors.New("failed to scan JSONMap: type assertion to []uint8 failed")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
