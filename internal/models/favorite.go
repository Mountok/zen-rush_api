package models

type Favorite struct {
	UserID     uint `gorm:"primaryKey" json:"user_id"`
	ActivityID uint `gorm:"primaryKey" json:"activity_id"`
}
