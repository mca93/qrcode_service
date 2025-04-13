package validators

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mca93/qrcode_service/models"
)

func ValidateTemplateCreate(req models.TemplateCreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.ClientAppID) == "" {
		return errors.New("clientAppId is required")
	}
	return validateStyle(req.Style)
}

func ValidateTemplateUpdate(req models.TemplateUpdateRequest) error {
	if req.Style != nil {
		return validateStyle(*req.Style)
	}
	return nil
}

func validateStyle(style models.QRCodeStyle) error {
	// Validate size
	if style.Size <= 0 {
		return errors.New("invalid size")
	}

	// Validate margin
	if style.Margin < 0 {
		return errors.New("margin cannot be negative")
	}

	// Validate error correction level
	switch style.ErrorCorrection {
	case models.ErrorCorrectionL, models.ErrorCorrectionM, models.ErrorCorrectionQ, models.ErrorCorrectionH:
	default:
		return errors.New("invalid errorCorrection (must be L, M, Q, or H)")
	}

	// Validate HEX color fields
	if !isValidHexColor(style.ForegroundColor) {
		return errors.New("invalid foregroundColor (must be a valid HEX color)")
	}
	if !isValidHexColor(style.BackgroundColor) {
		return errors.New("invalid backgroundColor (must be a valid HEX color)")
	}
	if style.Gradient && !isValidHexColor(style.GradientColor) {
		return errors.New("invalid gradientColor (must be a valid HEX color)")
	}
	if !isValidHexColor(style.BorderColor) {
		return errors.New("invalid borderColor (must be a valid HEX color)")
	}

	return nil
}

// isValidHexColor validates if a string is a valid HEX color.
func isValidHexColor(color string) bool {
	hexColorRegex := regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)
	return hexColorRegex.MatchString(color)
}
