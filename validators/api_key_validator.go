package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mca93/qrcode_service/models"
)

var oneWordRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func ValidateApiKeyCreate(req models.ApiKeyCreateRequest, clientAppID string) (string, error) {
	// Validação do nome (uma palavra)
	if strings.TrimSpace(req.Name) == "" {
		return "", errors.New("name is required")
	}

	if !oneWordRegex.MatchString(req.Name) {
		return "", errors.New("name must be a single word (alphanumeric only)")
	}

	// Validação do status
	if !isValidApiKeyStatus(req.Status) {
		return "", errors.New("invalid API key status")
	}

	// Geração do prefixo
	timestamp := time.Now().UTC().Format("20060102150405")
	keyPrefix := fmt.Sprintf("%s_%s_%s", req.Name, clientAppID, timestamp)

	return keyPrefix, nil
}

func isValidApiKeyStatus(status models.ApiKeyStatus) bool {
	switch status {
	case models.ApiKeyStatusUnspecified, models.ApiKeyStatusActive, models.ApiKeyStatusRevoked:
		return true
	default:
		return false
	}
}
