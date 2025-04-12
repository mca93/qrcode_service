package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mca93/qrcode_service/config"
	"github.com/mca93/qrcode_service/models"
	"github.com/mca93/qrcode_service/validators"

	"github.com/gin-gonic/gin"
)

// POST /v1/clientapps
func CreateClientApp(c *gin.Context) {
	var req models.ClientAppCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := validators.ValidateClientAppCreate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientApp := models.ClientApp{
		ID:           uuid.NewString(),
		Name:         req.Name,
		ContactEmail: req.ContactEmail,
		Status:       models.ClientAppStatusActive,
		CreatedAt:    time.Now(),
	}

	if err := config.DB.Create(&clientApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create client app"})
		return
	}

	c.JSON(http.StatusOK, clientApp)
}

// GET /v1/clientapps
func ListClientApps(c *gin.Context) {
	// Query Params
	status := c.DefaultQuery("status", string(models.ClientAppStatusUnspecified))
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// Conversão
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Consulta filtrada e paginada
	var apps []models.ClientApp
	var total int64

	query := config.DB.Model(&models.ClientApp{}).Where("status = ?", status)
	query.Count(&total)

	err := query.Offset(offset).Limit(pageSize).Find(&apps).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar client apps"})
		return
	}

	// Monta resposta
	var responses []models.ClientAppResponse
	for _, a := range apps {
		responses = append(responses, models.ClientAppResponse{
			ID:           a.ID,
			Name:         a.Name,
			ContactEmail: a.ContactEmail,
			Status:       a.Status,
			CreatedAt:    a.CreatedAt,
			UpdatedAt:    a.CreatedAt,
			DeletedAt:    nil,
		})
	}

	c.JSON(http.StatusOK, models.ClientAppListResponse{
		ClientApps: responses,
		TotalCount: int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
		HasNext:    offset+pageSize < int(total),
		HasPrev:    page > 1,
		NextPage:   page + 1,
		PrevPage:   page - 1,
		FirstPage:  1,
		LastPage:   (int(total) + pageSize - 1) / pageSize,
		FirstItem:  offset + 1,
		LastItem:   offset + len(apps),
		ItemsCount: len(apps),
		Items:      responses,
	})
}

// GET /v1/clientapps/:id
func GetClientApp(c *gin.Context) {
	clientAppID := c.Param("id")
	var clientApp models.ClientApp

	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client app not found"})
		return
	}

	c.JSON(http.StatusOK, clientApp)
}

// PUT /v1/clientapps/:id
func UpdateClientApp(c *gin.Context) {
	clientAppID := c.Param("id")
	var req models.ClientAppUpdateRequest

	// Bind JSON request to the update request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate the update request
	if err := validators.ValidateClientAppUpdate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the existing client app from the database
	var clientApp models.ClientApp
	if err := config.DB.First(&clientApp, "id = ?", clientAppID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client app not found"})
		return
	}

	// Update fields if provided in the request
	if req.Name != "" {
		clientApp.Name = req.Name
	}
	if req.ContactEmail != "" {
		clientApp.ContactEmail = req.ContactEmail
	}
	if req.Status != "" {
		clientApp.Status = req.Status
	}

	// Save the updated client app to the database
	if err := config.DB.Save(&clientApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update client app"})
		return
	}

	// Return the updated client app as a response
	c.JSON(http.StatusOK, models.ClientAppResponse{
		ID:           clientApp.ID,
		Name:         clientApp.Name,
		ContactEmail: clientApp.ContactEmail,
		Status:       clientApp.Status,
		CreatedAt:    clientApp.CreatedAt,
		UpdatedAt:    clientApp.UpdatedAt,
	})
}
