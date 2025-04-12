package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/routes"

	cmd "github.com/mca93/qrcode_service/cmd/swagger"
	"github.com/mca93/qrcode_service/config"
)

func main() {
	r := gin.Default() // cria um novo router
	cmd.RegisterSwagger(r)
	config.InitDB()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
