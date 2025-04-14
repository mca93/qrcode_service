package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type QRCodeErrorCorrection string

const (
	ErrorCorrectionL QRCodeErrorCorrection = "L"
	ErrorCorrectionM QRCodeErrorCorrection = "M"
	ErrorCorrectionQ QRCodeErrorCorrection = "Q"
	ErrorCorrectionH QRCodeErrorCorrection = "H"
)

type QRCodeStyle struct {
	Shape           string                `json:"shape"`
	ForegroundColor string                `json:"foregroundColor"`
	BackgroundColor string                `json:"backgroundColor"`
	Size            int                   `json:"size"`
	Margin          int                   `json:"margin"`
	CornerRadius    int                   `json:"cornerRadius"`
	Gradient        bool                  `json:"gradient"`
	GradientColor   string                `json:"gradientColor"`
	GradientAngle   int                   `json:"gradientAngle"`
	Border          int                   `json:"border"`
	BorderColor     string                `json:"borderColor"`
	LogoURL         string                `json:"logoUrl"`
	ErrorCorrection QRCodeErrorCorrection `json:"errorCorrection"`
}

type Template struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ClientAppID string    `gorm:"not null" json:"clientAppId"`                   // Foreign key to ClientApp
	ClientApp   ClientApp `gorm:"foreignKey:ClientAppID;references:ID" json:"-"` // Association with ClientApp
	Style       JSONMap   `gorm:"type:json" json:"style"`                        // JSON field
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TemplateCreateRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	ClientAppID string      `json:"clientAppId" binding:"required"`
	Style       QRCodeStyle `json:"style" binding:"required"`
}

type TemplateUpdateRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Style       *QRCodeStyle `json:"style"`
}

// JSONMap is a custom type to handle JSON fields in the database.
type JSONMap map[string]interface{}

// Scan implements the sql.Scanner interface for JSONMap.
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]uint8)
	if !ok {
		return errors.New("failed to scan JSONMap: type assertion to []uint8 failed")
	}

	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for JSONMap.
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
