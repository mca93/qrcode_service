package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mca93/qrcode_service/models"
)

const (
	maxNameLength        = 100
	maxDescriptionLength = 500
	minQRSize            = 100
	maxQRSize            = 1000
	colorRegex           = "^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
)

// ValidateTemplateCreate validates a TemplateCreateRequest
func ValidateTemplateCreate(req models.TemplateCreateRequest) error {
	if err := validateName(req.Name); err != nil {
		return err
	}

	if err := validateDescription(req.Description); err != nil {
		return err
	}

	if err := validateClientAppID(req.ClientAppID); err != nil {
		return err
	}

	if err := validateDefinition(req.Definition); err != nil {
		return err
	}

	if err := validateShape(req.Shape); err != nil {
		return err
	}

	if err := validateColor(req.ForegroundColor, "foregroundColor"); err != nil {
		return err
	}

	if err := validateColor(req.BackgroundColor, "backgroundColor"); err != nil {
		return err
	}

	if err := validateSize(req.Size); err != nil {
		return err
	}

	if err := validateLogoURL(req.LogoURL); err != nil {
		return err
	}

	if err := validateErrorCorrection(req.ErrorCorrection); err != nil {
		return err
	}

	return nil
}

// ValidateTemplateUpdate validates a TemplateUpdateRequest
func ValidateTemplateUpdate(req models.TemplateUpdateRequest) error {
	// Reuse the same validation as create since the fields are the same
	return ValidateTemplateCreate(models.TemplateCreateRequest(req))
}

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if len(name) > maxNameLength {
		return fmt.Errorf("name must be less than %d characters", maxNameLength)
	}

	return nil
}

func validateDescription(description string) error {
	if len(description) > maxDescriptionLength {
		return fmt.Errorf("description must be less than %d characters", maxDescriptionLength)
	}
	return nil
}

func validateClientAppID(clientAppID string) error {
	if strings.TrimSpace(clientAppID) == "" {
		return errors.New("clientAppId is required")
	}
	// Add more validation if needed (like UUID format check)
	return nil
}

func validateDefinition(definition models.Definition) error {
	if len(definition) == 0 {
		return errors.New("definition must contain at least one field")
	}

	fieldNames := make(map[string]bool)
	for _, field := range definition {
		if err := field.Validate(); err != nil {
			return fmt.Errorf("invalid field: %w", err)
		}

		// Check for duplicate field names
		if fieldNames[field.Name] {
			return fmt.Errorf("duplicate field name: %s", field.Name)
		}
		fieldNames[field.Name] = true

		// Validate field-specific validations
		if err := validateFieldValidations(field); err != nil {
			return fmt.Errorf("invalid validations for field %s: %w", field.Name, err)
		}
	}

	return nil
}

func validateFieldValidations(field models.Field) error {
	switch field.Type {
	case models.FieldTypeText:
		return validateTextValidations(field.Validations)
	case models.FieldTypeNumber:
		return validateNumberValidations(field.Validations)
	case models.FieldTypeMedia:
		return validateMediaValidations(field.Validations)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type)
	}
}

func validateTextValidations(validations map[string]interface{}) error {
	// Example validation for text fields
	if minLen, ok := validations["minLength"]; ok {
		if min, ok := minLen.(float64); ok {
			if min < 0 {
				return errors.New("minLength must be >= 0")
			}
		} else {
			return errors.New("minLength must be a number")
		}
	}

	if maxLen, ok := validations["maxLength"]; ok {
		if max, ok := maxLen.(float64); ok {
			if max <= 0 {
				return errors.New("maxLength must be > 0")
			}
		} else {
			return errors.New("maxLength must be a number")
		}
	}

	// Add more text validations as needed
	return nil
}

func validateNumberValidations(validations map[string]interface{}) error {
	// Example validation for number fields
	if minVal, ok := validations["min"]; ok {
		if _, ok := minVal.(float64); !ok {
			return errors.New("min must be a number, got %T")
		}
	}

	if maxVal, ok := validations["max"]; ok {
		if _, ok := maxVal.(float64); !ok {
			return errors.New("max must be a number")
		}
	}

	// Check if min <= max when both are present
	if minVal, ok1 := validations["min"]; ok1 {
		if maxVal, ok2 := validations["max"]; ok2 {
			if minVal.(float64) > maxVal.(float64) {
				return errors.New("min must be less than or equal to max")
			}
		}
	}

	return nil
}

func validateMediaValidations(validations map[string]interface{}) error {
	// Example validation for media fields
	if allowedTypes, ok := validations["allowedTypes"]; ok {
		if types, ok := allowedTypes.([]interface{}); ok {
			if len(types) == 0 {
				return errors.New("allowedTypes must contain at least one type")
			}
			for _, t := range types {
				if _, ok := t.(string); !ok {
					return errors.New("allowedTypes must be an array of strings")
				}
			}
		} else {
			return errors.New("allowedTypes must be an array")
		}
	}

	if maxSize, ok := validations["maxSize"]; ok {
		if size, ok := maxSize.(float64); ok {
			if size <= 0 {
				return errors.New("maxSize must be > 0")
			}
		} else {
			return errors.New("maxSize must be a number")
		}
	}

	return nil
}

func validateShape(shape string) error {
	if shape == "" {
		return nil // shape is optional
	}

	// Add shape validation if there are specific allowed shapes
	// Example:
	// allowedShapes := map[string]bool{"square": true, "circle": true, "rounded": true}
	// if !allowedShapes[shape] {
	//     return fmt.Errorf("invalid shape: %s. Allowed shapes are: square, circle, rounded", shape)
	// }

	return nil
}

func validateColor(color, fieldName string) error {
	if color == "" {
		return nil // color is optional
	}

	matched, err := regexp.MatchString(colorRegex, color)
	if err != nil {
		return fmt.Errorf("error validating %s: %v", fieldName, err)
	}
	if !matched {
		return fmt.Errorf("invalid %s format, must be a valid hex color (e.g. #RRGGBB or #RGB)", fieldName)
	}

	return nil
}

func validateSize(size int) error {
	if size < minQRSize || size > maxQRSize {
		return fmt.Errorf("size must be between %d and %d", minQRSize, maxQRSize)
	}
	return nil
}

func validateLogoURL(logoURL string) error {
	if logoURL == "" {
		return nil // logo is optional
	}

	// Add URL validation if needed
	// Example:
	// if _, err := url.ParseRequestURI(logoURL); err != nil {
	//     return fmt.Errorf("invalid logo URL: %v", err)
	// }

	return nil
}

func validateErrorCorrection(ec models.QRCodeErrorCorrection) error {
	switch ec {
	case models.ErrorCorrectionL, models.ErrorCorrectionM, models.ErrorCorrectionQ, models.ErrorCorrectionH, "":
		return nil
	default:
		return fmt.Errorf("invalid error correction level: %s. Valid options are L, M, Q, H", ec)
	}
}

// ValidateTemplateFilters validates template filter parameters
func ValidateTemplateFilters(active *bool, clientAppID string, createdAtFrom, createdAtTo *time.Time) error {
	if clientAppID != "" {
		if err := validateClientAppID(clientAppID); err != nil {
			return err
		}
	}

	if createdAtFrom != nil && createdAtTo != nil {
		if createdAtFrom.After(*createdAtTo) {
			return errors.New("createdAtFrom must be before or equal to createdAtTo")
		}
	}

	return nil
}
