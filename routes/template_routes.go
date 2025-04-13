package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mca93/qrcode_service/controllers"
)

func RegisterTemplateRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/templates")
	{
		v1.GET("", controllers.ListTemplates)
		v1.POST("", controllers.CreateTemplate)
		v1.GET("/:id", controllers.GetTemplate)
		v1.PATCH("/:id", controllers.UpdateTemplate)
		v1.POST("/:id/deactivate", controllers.DeactivateTemplate)
	}
}
