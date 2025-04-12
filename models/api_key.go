package models

import (
	"time"
)

type ApiKeyStatus string

const (
	ApiKeyStatusUnspecified ApiKeyStatus = "API_KEY_STATUS_UNSPECIFIED"
	ApiKeyStatusActive      ApiKeyStatus = "API_KEY_STATUS_ACTIVE"
	ApiKeyStatusRevoked     ApiKeyStatus = "API_KEY_STATUS_REVOKED"
)

type ApiKey struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	Name        string       `json:"name"`
	ClientAppID string       `json:"clientAppId"`
	KeyPrefix   string       `json:"keyPrefix"`
	Status      ApiKeyStatus `json:"status"`
	CreatedAt   time.Time    `json:"createdAt"`
	RevokedAt   *time.Time   `json:"revokedAt,omitempty"`
}
type ApiKeyCreateRequest struct {
	Name   string       `json:"name" binding:"required"` // Prefixo da chave
	Status ApiKeyStatus `json:"status"`                  // Status da chave
}
type ApiKeyUpdateRequest struct {
	Name        string       `json:"keyPrefix"`   // Prefixo da chave
	ClientAppID string       `json:"clientAppId"` // ID do aplicativo cliente
	Status      ApiKeyStatus `json:"status"`      // Status da chave
}
type ApiKeyResponse struct {
	ID          string       `json:"id"`                  // ID da chave
	ClientAppID string       `json:"clientAppId"`         // ID do aplicativo cliente
	Name        string       `json:"name"`                // Nome da chave
	KeyPrefix   string       `json:"keyPrefix"`           // Prefixo da chave
	Status      ApiKeyStatus `json:"status"`              // Status da chave
	CreatedAt   time.Time    `json:"createdAt"`           // Data de criação da chave
	RevokedAt   *time.Time   `json:"revokedAt,omitempty"` // Data de revogação da chave
}
type ApiKeyListResponse struct {
	ApiKeys    []ApiKeyResponse `json:"apiKeys"`    // Lista de chaves
	TotalCount int              `json:"totalCount"` // Contagem total de chaves
	Page       int              `json:"page"`       // Página atual
	PageSize   int              `json:"pageSize"`   // Tamanho da página
	TotalPages int              `json:"totalPages"` // Total de páginas
	HasNext    bool             `json:"hasNext"`    // Se há próxima página
	HasPrev    bool             `json:"hasPrev"`    // Se há página anterior
	NextPage   int              `json:"nextPage"`   // Próxima página
	PrevPage   int              `json:"prevPage"`   // Página anterior
	FirstPage  int              `json:"firstPage"`  // Primeira página
	LastPage   int              `json:"lastPage"`   // Última página
	FirstItem  int              `json:"firstItem"`  // Primeiro item
	LastItem   int              `json:"lastItem"`   // Último item
	ItemsCount int              `json:"itemsCount"` // Contagem de itens
	Items      []ApiKeyResponse `json:"items"`      // Itens da página
}
