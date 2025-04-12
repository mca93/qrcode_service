package models

import "time"

type QRCode struct {
    ID         string    `gorm:"primaryKey"`
    Type       string
    CreatedAt  time.Time
    ExpiresAt  *time.Time
    Status     string
    ScanCount  int64
    ImageURL   string
    DeepLinkURL string
    ClientAppID string
    TemplateID  string
}
