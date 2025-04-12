package models

type Template struct {
    ID          string `gorm:"primaryKey"`
    ClientAppID string
    Name        string
    Description string
    Active      bool
}
