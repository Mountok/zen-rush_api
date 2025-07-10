package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Activity struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:128" json:"name"`
	Description string         `json:"description"`
	Budget      int            `json:"budget"`
	Time        int            `json:"time"` // Сколько времени займёт (в часах)
	Weather     string         `gorm:"size:16" json:"weather"`
	Moods       pq.StringArray `gorm:"type:varchar(64)[]" json:"moods"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
