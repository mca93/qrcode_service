package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/controllers"
)

// RegisterTemplateRoutes registers all routes related to templates.
func RegisterTemplateRoutes(r *gin.Engine) {
	templateRoutes := r.Group("/v1/templates")
	{
		templateRoutes.GET("", controllers.NewTemplateController().ListTemplates)                      // List all templates
		templateRoutes.POST("", controllers.NewTemplateController().CreateTemplate)                    // Create a new template
		templateRoutes.GET("/:id", controllers.NewTemplateController().GetTemplate)                    // Get a specific template by ID
		templateRoutes.PATCH("/:id", controllers.NewTemplateController().UpdateTemplate)               // Update a specific template by ID
		templateRoutes.POST("/:id/deactivate", controllers.NewTemplateController().DeactivateTemplate) // Deactivate a specific template by ID
	}
}
