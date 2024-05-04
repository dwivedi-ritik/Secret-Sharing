package models

import (
	"time"
)

type User struct {
	Id         uint64     `gorm:"primaryKey"`
	Username   string     `gorm:"not null;unique"`
	Email      string     `gorm:"not null;unique"`
	Password   string     `gorm:"not null"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Messages   []Message  `gorm:"foreignKey:CreatedBy"`
	Encryption Encryption `gorm:"foreignKey:CreatedBy"`
}
