package validators

import (
	"errors"
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
