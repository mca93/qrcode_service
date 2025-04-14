package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type QRCodeErrorCorrection string

const (
	ErrorCorrectionL QRCodeErrorCorrection = "L"
	ErrorCorrectionM QRCodeErrorCorrection = "M"
	ErrorCorrectionQ QRCodeErrorCorrection = "Q"
	ErrorCorrectionH QRCodeErrorCorrection = "H"
)

// ---------- ENUMS ----------

type FieldType string

const (
	FieldTypeText   FieldType = "Text"
	FieldTypeNumber FieldType = "Number"
	FieldTypeMedia  FieldType = "Media"
)

func (ft FieldType) IsValid() bool {
	switch ft {
	case FieldTypeText, FieldTypeNumber, FieldTypeMedia:
		return true
	default:
		return false
	}
}

type MetadataType string

const (
	MetadataTypeContent MetadataType = "Content"
	MetadataTypeStyle   MetadataType = "Style"
)

func (mt MetadataType) IsValid() bool {
	switch mt {
	case MetadataTypeContent, MetadataTypeStyle:
		return true
	default:
		return false
	}
}

// ---------- STRUCTS ----------

type FieldDefinition struct {
	Name        string                 `json:"name"`
	Type        FieldType              `json:"type"`
	Validations map[string]interface{} `json:"validations"`
}

type MetadataDefinition struct {
	Type   MetadataType      `json:"type"`
	Fields []FieldDefinition `json:"fields"`
}

// ---------- DATABASE SERIALIZATION ----------

func (m MetadataDefinition) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MetadataDefinition) Scan(value interface{}) error {
	if value == nil {
		*m = MetadataDefinition{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan MetadataDefinition: expected []byte")
	}
	return json.Unmarshal(bytes, m)
}

// Optional: FieldDefinition validation logic
func (f FieldDefinition) Validate() error {
	if f.Name == "" {
		return errors.New("field name is required")
	}
	if !f.Type.IsValid() {
		return fmt.Errorf("invalid field type: %s", f.Type)
	}
	return nil
}

func (m MetadataDefinition) Validate() error {
	if !m.Type.IsValid() {
		return fmt.Errorf("invalid metadata type: %s", m.Type)
	}
	for _, field := range m.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

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

// --------- Template -----------

type Template struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ClientAppID string    `gorm:"not null" json:"clientAppId"`
	ClientApp   ClientApp `gorm:"foreignKey:ClientAppID;references:ID" json:"-"`
	// Style       JSONMap            `gorm:"type:json" json:"style"`
	Metadata  MetadataDefinition `gorm:"type:json" json:"metadata"`
	Active    bool               `json:"active"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

// --------- Requests -----------

type TemplateCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ClientAppID string `json:"clientAppId" binding:"required"`
	// Style       QRCodeStyle        `json:"style" binding:"required"`
	Metadata MetadataDefinition `json:"metadata"`
}

type TemplateUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Style       *QRCodeStyle        `json:"style"`
	Metadata *MetadataDefinition `json:"metadata"`
}
