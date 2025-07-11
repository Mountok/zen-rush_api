package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zenrush/backend/internal/db"
	"github.com/zenrush/backend/internal/models"
	"gorm.io/gorm/clause"
)

type MoodStatRequest struct {
	Mood string `json:"mood" binding:"required"`
	Date string `json:"date"` // YYYY-MM-DD, опционально
}

// POST /api/mood-stats
func SaveOrUpdateMoodStat(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req MoodStatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	var date time.Time
	var err error
	if req.Date != "" {
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
	} else {
		date = time.Now().Truncate(24 * time.Hour)
	}
	moodStat := models.MoodStat{
		UserID: userID,
		Date:   date,
		Mood:   req.Mood,
	}
	err = db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"mood"}),
	}).Create(&moodStat).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusCreated)
}

// GET /api/users/me/mood-stats?days=N
func GetMoodStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid days param"})
		return
	}
	fromDate := time.Now().AddDate(0, 0, -days+1).Truncate(24 * time.Hour)
	var stats []models.MoodStat
	err = db.DB.Where("user_id = ? AND date >= ?", userID, fromDate).Order("date asc").Find(&stats).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
