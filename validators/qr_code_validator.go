package validators

import (
	"errors"
	"fmt"
	"time"

	"github.com/mca93/qrcode_service/models"
)

// ValidateQRCodeCreate validates the QRCodeCreateRequest.
func ValidateQRCodeCreate(req models.QRCodeCreateRequest, clientAppID string) error {
	// Validate Type
	if req.Type != models.QRCodeTypeStable && req.Type != models.QRCodeTypeDynamic {
		return errors.New("invalid type: must be STABLE or DYNAMIC")
	}

	// Validate TemplateID
	if req.TemplateID == "" {
		return errors.New("templateId is required")
	}

	// Validate ClientAppID from the header
	if clientAppID == "" {
		return errors.New("clientAppId is required in the header")
	}

	if req.ThirdPartyRef == "" {
		return errors.New("thirdPartyRef is required")
	}
	// Ensure the ClientAppID in the header matches the request
	// if req.ClientAppID != "" && req.ClientAppID != clientAppID {
	// 	return errors.New("clientAppId in the request does not match the clientAppId in the header")
	// }

	// Validate ExpiresAt (if provided)
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		return errors.New("expiresAt cannot be in the past")
	}

	// Validate Data (if provided)
	if req.Data != nil {
		if err := validateCustomData(req.Data); err != nil {
			return err
		}
	}

	return nil
}

// ValidateQRCodeUpdate validates the QRCodeUpdateRequest.
func ValidateQRCodeUpdate(req models.QRCodeUpdateRequest) error {
	// Validate Type (if provided)
	if req.Type != "" && req.Type != models.QRCodeTypeStable && req.Type != models.QRCodeTypeDynamic {
		return errors.New("invalid type: must be STABLE or DYNAMIC")
	}

	// Validate ExpiresAt (if provided)
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		return errors.New("expiresAt cannot be in the past")
	}

	// Validate Data (if provided)
	if req.Data != nil {
		if err := validateCustomData(req.Data); err != nil {
			return err
		}
	}

	return nil
}

// ValidateQRCodeOwnership ensures the QR code belongs to the requesting ClientAppID.
func ValidateQRCodeOwnership(clientAppID string, qrCode models.QRCode) error {
	if clientAppID == "" {
		return errors.New("clientAppId is required in the header")
	}

	if qrCode.ClientAppID != clientAppID {
		return errors.New("you do not have permission to access this QR code")
	}
	return nil
}

// validateCustomData validates the custom data field.
func validateCustomData(data map[string]interface{}) error {
	// Example validation: Ensure no empty keys or values
	for key, value := range data {
		if key == "" {
			return errors.New("data contains an empty key")
		}
		if value == nil {
			return errors.New("data contains a nil value for key: " + key)
		}
	}
	return nil
}

// ValidateQRCodeData validates the Data field of a QR code against the template's Definition.
func ValidateQRCodeData(data models.JSONMap, template models.Template) error {
	// Iterate through the Definition to validate fields
	for _, field := range template.Definition {
		// Check if the field is marked as required
		if required, ok := field.Validations["required"].(bool); ok && required {
			// Ensure the required field exists in the Data map
			if _, exists := data[field.Name]; !exists {
				return fmt.Errorf("missing required field in data: %s", field.Name)
			}
		}

		// Additional validation for specific field types
		if err := validateFieldData(field, data[field.Name]); err != nil {
			return fmt.Errorf("validation failed for field '%s': %w", field.Name, err)
		}
	}

	// Validate that all keys in the Data field match the fields in the Definition
	validKeys := make(map[string]bool)
	for _, field := range template.Definition {
		validKeys[field.Name] = true
	}

	for key := range data {
		if !validKeys[key] {
			return fmt.Errorf("invalid key in data: %s", key)
		}
	}

	return nil
}

// validateFieldData validates the value of a specific field based on its type and validations.
func validateFieldData(field models.Field, value interface{}) error {
	// Skip validation if the value is nil
	if value == nil {
		return nil
	}

	switch field.Type {
	case models.FieldTypeText:
		if maxLength, ok := field.Validations["maxLength"].(float64); ok {
			if len(fmt.Sprintf("%v", value)) > int(maxLength) {
				return fmt.Errorf("value exceeds maxLength of %d", int(maxLength))
			}
		}
	case models.FieldTypeNumber:
		if min, ok := field.Validations["min"].(float64); ok {
			if num, ok := value.(float64); ok && num < min {
				return fmt.Errorf("value is less than min of %f", min)
			}
		}
		if max, ok := field.Validations["max"].(float64); ok {
			if num, ok := value.(float64); ok && num > max {
				return fmt.Errorf("value exceeds max of %f", max)
			}
		}
	case models.FieldTypeMedia:
		// Add specific validations for media fields if needed
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type)
	}

	return nil
}
