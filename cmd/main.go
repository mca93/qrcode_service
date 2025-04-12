package main

import (
	"github.com/mca93/qrcode_service/routes"

	"github.com/mca93/qrcode_service/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // cria um novo router
	RegisterSwagger(r)
	config.InitDB()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
