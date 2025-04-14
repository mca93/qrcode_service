package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// Regista as rotas do ClientApp
	RegisterClientAppRoutes(router)

	// Regista as rotas do Template
	RegisterTemplateRoutes(router)

	// Regista as rotas do QRCode
	RegisterQRCodeRoutes(router)

}
