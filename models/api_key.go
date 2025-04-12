package models

import "time"

type ApiKey struct {
    ID          string    `gorm:"primaryKey"`
    ClientAppID string
    KeyPrefix   string
    Status      string
    CreatedAt   time.Time
    RevokedAt   *time.Time
}
