package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zenrush/backend/internal/db"
	"github.com/zenrush/backend/internal/models"
)

// Получить все избранные активности пользователя
func ListFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	var favorites []models.Favorite
	if err := db.DB.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	var activityIDs []uint
	for _, f := range favorites {
		activityIDs = append(activityIDs, f.ActivityID)
	}
	var activities []models.Activity
	if len(activityIDs) > 0 {
		db.DB.Where("id IN ?", activityIDs).Find(&activities)
	}
	c.JSON(http.StatusOK, activities)
}

// Добавить активность в избранное
func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	activityID, err := strconv.Atoi(c.Param("activity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity_id"})
		return
	}
	fav := models.Favorite{UserID: userID, ActivityID: uint(activityID)}
	if err := db.DB.Create(&fav).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error or already exists"})
		return
	}
	c.Status(http.StatusCreated)
}

// Удалить активность из избранного
func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	activityID, err := strconv.Atoi(c.Param("activity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity_id"})
		return
	}
	if err := db.DB.Delete(&models.Favorite{}, "user_id = ? AND activity_id = ?", userID, activityID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}
