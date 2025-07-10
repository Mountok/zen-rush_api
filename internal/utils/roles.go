package utils

import "github.com/gin-gonic/gin"

func IsModeratorOrAdmin(c *gin.Context) bool {
	role, ok := c.Get("role")
	if !ok {
		return false
	}
	r, ok := role.(string)
	if !ok {
		return false
	}
	return r == "admin" || r == "moderator"
}
