package models

import "time"

type History struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	ActivityID uint      `gorm:"index" json:"activity_id"`
	ViewedAt   time.Time `gorm:"autoCreateTime" json:"viewed_at"`
}
