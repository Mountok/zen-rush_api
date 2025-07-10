package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null;size:64" json:"username"`
	PasswordHash string    `gorm:"not null;size:128" json:"-"`
	Role         string    `gorm:"type:varchar(16);default:user" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
