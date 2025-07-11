package models

import "time"

type MoodStat struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserID uint      `gorm:"not null;index" json:"user_id"`
	Date   time.Time `gorm:"type:date;not null;index" json:"date"`
	Mood   string    `gorm:"type:varchar(64);not null" json:"mood"`
}
