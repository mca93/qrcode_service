package validators

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mca93/qrcode_service/models"
)

// ValidateTemplateCreate validates the TemplateCreateRequest.
func ValidateTemplateCreate(req models.TemplateCreateRequest) error {
	// Validate Name
	if req.Name == "" {
		return errors.New("name is required")
	}

	// Validate ClientAppID
	if req.ClientAppID == "" {
		return errors.New("clientAppId is required")
	}

	// Validate Metadata
	if len(req.Metadata) != 2 {
		return errors.New("metadata must contain exactly two objects")
	}
	for _, metadata := range req.Metadata {
		if err := validateMetadata(metadata); err != nil {
			return fmt.Errorf("invalid metadata: %w", err)
		}
	}

	return nil
}

// ValidateTemplateUpdate validates the TemplateUpdateRequest.
func ValidateTemplateUpdate(req models.TemplateUpdateRequest) error {
	// Validate Metadata (if provided)
	if req.Metadata != nil {
		if len(*req.Metadata) != 2 {
			return errors.New("metadata must contain exactly two objects")
		}
		for _, metadata := range *req.Metadata {
			if err := validateMetadata(metadata); err != nil {
				return fmt.Errorf("invalid metadata: %w", err)
			}
		}
	}

	return nil
}

// validateMetadata validates a single MetadataDefinition.
func validateMetadata(metadata models.MetadataDefinition) error {
	// Validate Metadata Type
	if !metadata.Type.IsValid() {
		return fmt.Errorf("invalid metadata type: %s", metadata.Type)
	}

	// Validate Fields
	for _, field := range metadata.Fields {
		if metadata.Type == models.MetadataTypeStyle {
			// Validate the Style object
			if err := validateQRCodeStyle(field.Validations); err != nil {
				return fmt.Errorf("invalid style object: %w", err)
			}
		} else {
			// Validate other field types
			if err := field.Validate(); err != nil {
				return fmt.Errorf("invalid field: %w", err)
			}
		}
	}

	return nil
}

// validateQRCodeStyle validates the QRCodeStyle object.
func validateQRCodeStyle(validations map[string]interface{}) error {
	// Convert the map to a QRCodeStyle struct
	var style models.QRCodeStyle
	data, err := json.Marshal(validations)
	if err != nil {
		return fmt.Errorf("failed to marshal style validations: %w", err)
	}
	if err := json.Unmarshal(data, &style); err != nil {
		return fmt.Errorf("failed to unmarshal style validations: %w", err)
	}

	// Validate individual fields in QRCodeStyle
	if style.Shape == "" {
		return errors.New("shape is required")
	}
	if style.ForegroundColor == "" {
		return errors.New("foregroundColor is required")
	}
	if style.BackgroundColor == "" {
		return errors.New("backgroundColor is required")
	}
	if style.Size <= 0 {
		return errors.New("size must be greater than 0")
	}
	if style.Margin < 0 {
		return errors.New("margin cannot be negative")
	}
	if style.CornerRadius < 0 {
		return errors.New("cornerRadius cannot be negative")
	}
	if style.Gradient && style.GradientColor == "" {
		return errors.New("gradientColor is required when gradient is true")
	}
	if style.Border < 0 {
		return errors.New("border cannot be negative")
	}
	if style.BorderColor == "" && style.Border > 0 {
		return errors.New("borderColor is required when border is greater than 0")
	}
	if style.ErrorCorrection != models.ErrorCorrectionL &&
		style.ErrorCorrection != models.ErrorCorrectionM &&
		style.ErrorCorrection != models.ErrorCorrectionQ &&
		style.ErrorCorrection != models.ErrorCorrectionH {
		return errors.New("invalid errorCorrection value")
	}

	return nil
}
