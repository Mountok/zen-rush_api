package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zenrush/backend/internal/db"
	"github.com/zenrush/backend/internal/models"
)

// Получить последние 10 просмотренных активностей пользователя
func ListHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	var history []models.History
	if err := db.DB.Where("user_id = ?", userID).Order("viewed_at desc").Limit(10).Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	var activityIDs []uint
	for _, h := range history {
		activityIDs = append(activityIDs, h.ActivityID)
	}
	var activities []models.Activity
	if len(activityIDs) > 0 {
		db.DB.Where("id IN ?", activityIDs).Find(&activities)
	}
	c.JSON(http.StatusOK, activities)
}

// Добавить просмотр активности в историю
func AddHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	activityID, err := strconv.Atoi(c.Param("activity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity_id"})
		return
	}
	h := models.History{UserID: userID, ActivityID: uint(activityID)}
	if err := db.DB.Create(&h).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusCreated)
}
