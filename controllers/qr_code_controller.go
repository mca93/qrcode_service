package controllers

import (
	"net/http"

	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"

	"image/color"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func ListQRCodes(c *gin.Context) {
	var codes []models.QRCode
	config.DB.Find(&codes)
	c.JSON(http.StatusOK, gin.H{"qr_codes": codes})
}

func CreateQRCode(c *gin.Context) {
	var input models.QRCode
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}

func GetQRCodeImage(c *gin.Context) {
	id := c.Param("id")
	var qr models.QRCode

	if err := config.DB.First(&qr, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR Code not found"})
		return
	}

	png, err := qrcode.Encode(qr.DeepLinkURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

// Função para QRCode com estilo
func GenerateStyledQRCode(url string) ([]byte, error) {
	qr, err := qrcode.New(url, qrcode.High)
	if err != nil {
		return nil, err
	}
	qr.BackgroundColor = color.White
	qr.ForegroundColor = color.Black
	return qr.PNG(256)
}
