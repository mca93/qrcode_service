package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"os"

	"github.com/mca93/qrcode_service/models"
	qrcode "github.com/yeqown/go-qrcode/v2"
	qs "github.com/yeqown/go-qrcode/writer/standard"
)

// QRCodeService handles QR code generation for a specific QRCode instance.
type QRCodeService struct {
	QRCode   *models.QRCode
	Template *models.Template
}

// NewQRCodeService initializes a new QRCodeService with the given QRCode and Template objects.
func NewQRCodeService(qr *models.QRCode, template *models.Template) *QRCodeService {
	return &QRCodeService{
		QRCode:   qr,
		Template: template,
	}
}

// GenerateBase64Image generates a base64-encoded PNG image of the QR code.
func (s *QRCodeService) GenerateBase64Image() (string, error) {
	dataToEncode := s.getDataToEncode()

	_, imageOptions, err := s.getQRCodeOptions()
	if err != nil {
		return "", fmt.Errorf("failed to get QR code options: %w", err)
	}

	// qrCode, err := qrcode.New(dataToEncode)
	qrCode, err := qrcode.NewWith(dataToEncode, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest))
	if err != nil {
		return "", err
	}

	imageBytes, err := s.generateQRCodeImage(qrCode, imageOptions)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageBytes), nil
}

// getDataToEncode returns the data to be encoded in the QR code.
func (s *QRCodeService) getDataToEncode() string {
	if s.QRCode.DeepLinkURL != "" {
		return s.QRCode.DeepLinkURL
	}
	return fmt.Sprintf("https://yourdomain.com/qrcode/%s", s.QRCode.ID)
}

type Option interface{}

// getQRCodeOptions initializes QR code options based on the style from the Template model.
func (s *QRCodeService) getQRCodeOptions() ([]qrcode.EncodeOption, []qs.ImageOption, error) {
	// Retrieve the style object from the metadata
	style, err := s.getStyleMetadata()
	if err != nil {
		return nil, nil, err
	}

	imgOptions := []qs.ImageOption{
		// qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
		qs.WithBorderWidth(8),
		qs.WithFgColorRGBHex(style["foregroundColor"].(string)),
		qs.WithBgColorRGBHex(style["backgroundColor"].(string)),
	}

	qrCodeOptions := []qrcode.EncodeOption{
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	}

	if err := s.applyShapeOption(style, &imgOptions); err != nil {
		return nil, nil, err
	}

	if err := s.applyLogoOption(style, &imgOptions); err != nil {
		return nil, nil, err
	}

	return qrCodeOptions, imgOptions, nil
}

// getStyleMetadata retrieves the metadata object of type "Style" from the Template.
func (s *QRCodeService) getStyleMetadata() (map[string]interface{}, error) {
	for _, metadata := range s.Template.Metadata {
		fmt.Printf("Fetching style metadata from template ID: %+v\n", s.Template)

		if metadata.Type == "Style" {
			// Convert the fields into a map for easier access
			style := make(map[string]interface{})
			for name, field := range metadata.Fields[0].Validations {

				style[name] = field
			}
			return style, nil
		}
	}
	return nil, fmt.Errorf("style metadata not found in template")
}

// applyShapeOption applies the shape option if specified in the style.
func (s *QRCodeService) applyShapeOption(style map[string]interface{}, options *[]qs.ImageOption) error {
	if shape, ok := style["shape"].(string); ok {
		switch shape {
		case "circle":
			*options = append(*options, qs.WithCircleShape())
		case "square":
			// *options = append(*options, qs.WithCustomShape())
		default:
			return fmt.Errorf("unsupported shape: %s", shape)
		}
	}
	return nil
}

// applyLogoOption applies the logo option if specified in the style.
func (s *QRCodeService) applyLogoOption(style map[string]interface{}, options *[]qs.ImageOption) error {
	if logoURL, ok := style["logoUrl"].(string); ok && logoURL != "" {
		logoBytes, err := os.ReadFile(logoURL)
		if err != nil {
			return err
		}
		img, _, err := image.Decode(bytes.NewBuffer(logoBytes))
		if err != nil {
			return fmt.Errorf("failed to read logo file: %w", err)
		}
		opt := qs.WithLogoImage(img)
		*options = append(*options, opt)
	}
	return nil
}

// // applyOptions applies all options to the QR code.
// func applyOptions(qrCode *qrcode.QRCode, options []qrcode.EncodeOption) error {
// 	for _, option := range options {
// 		if err := qrCode.Apply(option); err != nil {
// 			return fmt.Errorf("failed to apply QR code option: %w", err)
// 		}
// 	}
// 	return nil
// }

type CustomWriteCloser interface {
	Close() error
	Write(p []byte) error
}

type customWriteCloser struct {
	buffer *bytes.Buffer
}

func (c *customWriteCloser) Close() error {
	return nil
}
func (c *customWriteCloser) Write(p []byte) (int, error) {
	return c.buffer.Write(p)
}

// generateQRCodeImage generates the QR code image and returns the bytes.
func (s *QRCodeService) generateQRCodeImage(qrCode *qrcode.QRCode, imageOptions []qs.ImageOption) ([]byte, error) {
	buf := new(bytes.Buffer)
	customWriteCloser := &customWriteCloser{buffer: buf}
	writer := qs.NewWithWriter(customWriteCloser, imageOptions...)

	// Finally save with both configurations
	if err := qrCode.Save(writer); err != nil {
		return nil, fmt.Errorf("failed to save QR code: %w", err)
	}

	return buf.Bytes(), nil
}
