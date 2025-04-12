package models

import (
	"time"
)

type ClientAppStatus string

const (
	ClientAppStatusUnspecified ClientAppStatus = "CLIENT_APP_STATUS_UNSPECIFIED"
	ClientAppStatusActive      ClientAppStatus = "CLIENT_APP_STATUS_ACTIVE"
	ClientAppStatusSuspended   ClientAppStatus = "CLIENT_APP_STATUS_SUSPENDED"
	ClientAppStatusDeleted     ClientAppStatus = "CLIENT_APP_STATUS_DELETED"
)

type ClientApp struct {
	ID           string          `gorm:"primaryKey" json:"ID"`
	Name         string          `json:"Name"`
	ContactEmail string          `json:"ContactEmail"`
	Status       ClientAppStatus `gorm:"default:CLIENT_APP_STATUS_ACTIVE" json:"Status"`
	CreatedAt    time.Time       `json:"CreatedAt"`
	UpdatedAt    time.Time       `json:"UpdatedAt"` // Data de atualização do aplicativo

}
type ClientAppCreateRequest struct {
	Name         string          `json:"Name"`         // Nome do aplicativo
	ContactEmail string          `json:"ContactEmail"` // Email de contato do cliente
	Status       ClientAppStatus `json:"Status"`       // Status do aplicativo
}
type ClientAppUpdateRequest struct {
	Name         string          `json:"Name"`         // Nome do aplicativo
	ContactEmail string          `json:"ContactEmail"` // Email de contato do cliente
	Status       ClientAppStatus `json:"Status"`       // Status do aplicativo
}
type ClientAppResponse struct {
	ID           string          `json:"ID"`           // ID do aplicativo
	Name         string          `json:"Name"`         // Nome do aplicativo
	ContactEmail string          `json:"ContactEmail"` // Email de contato do cliente
	Status       ClientAppStatus `json:"Status"`       // Status do aplicativo
	CreatedAt    time.Time       `json:"CreatedAt"`    // Data de criação do aplicativo
	UpdatedAt    time.Time       `json:"UpdatedAt"`    // Data de atualização do aplicativo
	DeletedAt    *time.Time      `json:"DeletedAt"`    // Data de exclusão do aplicativo
}
type ClientAppListResponse struct {
	ClientApps []ClientAppResponse `json:"ClientApps"` // Lista de aplicativos
	TotalCount int                 `json:"TotalCount"` // Contagem total de aplicativos
	Page       int                 `json:"Page"`       // Página atual
	PageSize   int                 `json:"PageSize"`   // Tamanho da página
	TotalPages int                 `json:"TotalPages"` // Total de páginas
	HasNext    bool                `json:"HasNext"`    // Se há próxima página
	HasPrev    bool                `json:"HasPrev"`    // Se há página anterior
	NextPage   int                 `json:"NextPage"`   // Próxima página
	PrevPage   int                 `json:"PrevPage"`   // Página anterior
	FirstPage  int                 `json:"FirstPage"`  // Primeira página
	LastPage   int                 `json:"LastPage"`   // Última página
	FirstItem  int                 `json:"FirstItem"`  // Primeiro item
	LastItem   int                 `json:"LastItem"`   // Último item
	ItemsCount int                 `json:"ItemsCount"` // Contagem de itens
	Items      []ClientAppResponse `json:"Items"`      // Itens da página
}
