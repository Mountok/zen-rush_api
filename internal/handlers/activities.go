package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zenrush/backend/internal/db"
	"github.com/zenrush/backend/internal/models"
	"github.com/zenrush/backend/internal/utils"
	"gorm.io/gorm/clause"
)

// Получить список всех активностей (с фильтрами)
func ListActivities(c *gin.Context) {
	var activities []models.Activity
	q := db.DB.Model(&models.Activity{})

	if minBudget := c.Query("min_budget"); minBudget != "" {
		if v, err := strconv.Atoi(minBudget); err == nil {
			q = q.Where("budget >= ?", v)
		}
	}
	if maxBudget := c.Query("max_budget"); maxBudget != "" {
		if v, err := strconv.Atoi(maxBudget); err == nil {
			q = q.Where("budget <= ?", v)
		}
	}
	if timeParam := c.Query("time"); timeParam != "" {
		if v, err := strconv.Atoi(timeParam); err == nil {
			q = q.Where("time = ?", v)
		}
	}
	if mood := c.Query("mood"); mood != "" {
		q = q.Where("? = ANY(moods)", mood)
		// --- Сохраняем настроение пользователя в статистику ---
		userID, exists := c.Get("user_id")
		if exists {
			// Upsert в mood_stats
			db.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
				DoUpdates: clause.AssignmentColumns([]string{"mood"}),
			}).Create(&models.MoodStat{
				UserID: userID.(uint),
				Date:   time.Now().Truncate(24 * time.Hour),
				Mood:   mood,
			})
		}
	}
	if weather := c.Query("weather"); weather != "" {
		q = q.Where("weather = ?", weather)
	}
	if peopleCount := c.Query("people_count"); peopleCount != "" {
		if v, err := strconv.Atoi(peopleCount); err == nil {
			q = q.Where("people_count = ?", v)
		}
	}

	if err := q.Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, activities)
}

// Получить одну активность по id
func GetActivity(c *gin.Context) {
	var activity models.Activity
	id := c.Param("id")
	if err := db.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, activity)
}

// Создать новую активность (только для модератора или админа)
func CreateActivity(c *gin.Context) {
	if !utils.IsModeratorOrAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req models.Activity
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// Обновить существующую активность (только для модератора или админа)
func UpdateActivity(c *gin.Context) {
	if !utils.IsModeratorOrAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	var activity models.Activity
	if err := db.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var req models.Activity
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	activity.Name = req.Name
	activity.Description = req.Description
	activity.Budget = req.Budget
	activity.Time = req.Time
	activity.Weather = req.Weather
	activity.PeopleCount = req.PeopleCount
	activity.Moods = req.Moods
	if err := db.DB.Save(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, activity)
}

// Удалить активность (только для модератора или админа)
func DeleteActivity(c *gin.Context) {
	if !utils.IsModeratorOrAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	if err := db.DB.Delete(&models.Activity{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}
