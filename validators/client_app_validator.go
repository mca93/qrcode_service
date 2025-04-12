package validators

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mca93/qrcode_service/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateClientAppCreate(req models.ClientAppCreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if !emailRegex.MatchString(req.ContactEmail) {
		return errors.New("invalid email format")
	}

	if req.Status == "" {
		req.Status = models.ClientAppStatusActive
	}

	return nil
}

func ValidateClientAppUpdate(req models.ClientAppUpdateRequest) error {
	if strings.TrimSpace(req.Name) == "" && req.ContactEmail == "" && req.Status == "" {
		return errors.New("at least one field (name, contact_email, status) is required")
	}

	if req.ContactEmail != "" {
		if !emailRegex.MatchString(req.ContactEmail) {
			return errors.New("invalid email format")
		}
	}
	if req.Status != "" {
		if !isValidStatus(req.Status) {
			return errors.New("invalid status value")
		}
	}

	return nil
}

func isValidStatus(status models.ClientAppStatus) bool {
	switch status {
	case models.ClientAppStatusUnspecified, models.ClientAppStatusActive, models.ClientAppStatusSuspended, models.ClientAppStatusDeleted:
		return true
	default:
		return false
	}
}
