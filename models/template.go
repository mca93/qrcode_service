package models

import (
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
	ID          string      `gorm:"primaryKey" json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ClientAppID string      `json:"clientAppId"`
	Style       QRCodeStyle `gorm:"type:jsonb" json:"style"`
	Active      bool        `json:"active"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
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
