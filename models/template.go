package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ---------- ENUMS ----------
type QRCodeErrorCorrection string

const (
	ErrorCorrectionL QRCodeErrorCorrection = "L"
	ErrorCorrectionM QRCodeErrorCorrection = "M"
	ErrorCorrectionQ QRCodeErrorCorrection = "Q"
	ErrorCorrectionH QRCodeErrorCorrection = "H"
)

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

// ---------- STRUCTS ----------

// Field represents a single field in the template definition.
type Field struct {
	Name        string                 `json:"name"`
	Type        FieldType              `json:"type"`
	Validations map[string]interface{} `json:"validations"`
}

// Definition is an array of fields that defines the structure of a template.
type Definition []Field

// ---------- DATABASE SERIALIZATION ----------

// Value implements the `driver.Valuer` interface for Definition.
func (d Definition) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the `sql.Scanner` interface for Definition.
func (d *Definition) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Definition: expected []byte")
	}

	return json.Unmarshal(bytes, d)
}

// ---------- VALIDATION ----------

// Validate validates a single field.
func (f Field) Validate() error {
	if f.Name == "" {
		return errors.New("field name is required")
	}
	if !f.Type.IsValid() {
		return fmt.Errorf("invalid field type: %s", f.Type)
	}
	return nil
}

// Validate validates the entire Definition.
func (d Definition) Validate() error {
	for _, field := range d {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ---------- TEMPLATE MODEL ----------

type Template struct {
	ID              string                `gorm:"primaryKey" json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	ClientAppID     string                `gorm:"not null" json:"clientAppId"`
	ClientApp       ClientApp             `gorm:"foreignKey:ClientAppID;references:ID" json:"-"`
	Definition      Definition            `gorm:"type:json" json:"definition"` // Updated to use Definition
	Shape           string                `json:"shape"`
	ForegroundColor string                `json:"foregroundColor"`
	BackgroundColor string                `json:"backgroundColor"`
	Size            int                   `json:"size"`
	LogoURL         string                `json:"logoUrl"`
	ErrorCorrection QRCodeErrorCorrection `json:"errorCorrection"`

	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ---------- REQUEST STRUCTS ----------

type TemplateCreateRequest struct {
	Name            string                `json:"name" binding:"required"`
	Description     string                `json:"description"`
	ClientAppID     string                `json:"clientAppId" binding:"required"`
	Definition      Definition            `json:"definition"` // Updated to use Definition
	Shape           string                `json:"shape"`
	ForegroundColor string                `json:"foregroundColor"`
	BackgroundColor string                `json:"backgroundColor"`
	Size            int                   `json:"size"`
	LogoURL         string                `json:"logoUrl"`
	ErrorCorrection QRCodeErrorCorrection `json:"errorCorrection"`
}

type TemplateUpdateRequest struct {
	Name            string                `json:"name" binding:"required"`
	Description     string                `json:"description"`
	ClientAppID     string                `json:"clientAppId" binding:"required"`
	Definition      Definition            `json:"definition"` // Updated to use Definition
	Shape           string                `json:"shape"`
	ForegroundColor string                `json:"foregroundColor"`
	BackgroundColor string                `json:"backgroundColor"`
	Size            int                   `json:"size"`
	LogoURL         string                `json:"logoUrl"`
	ErrorCorrection QRCodeErrorCorrection `json:"errorCorrection"`
}
