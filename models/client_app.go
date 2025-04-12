package models

import "time"

type ClientApp struct {
    ID           string    `gorm:"primaryKey"`
    Name         string
    ContactEmail string
    Status       string
    CreatedAt    time.Time
}
