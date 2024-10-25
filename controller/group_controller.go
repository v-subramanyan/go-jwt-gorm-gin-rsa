package controller

import (
	"jwt/initializers"
	"jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateGroup handles creating a new group
func CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DBConn.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

// GetGroup retrieves a single group by ID
func GetGroup(c *gin.Context) {
	var group models.Group
	id := c.Param("id")

	if err := initializers.DBConn.Preload("Members").First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

// UpdateGroup handles updating a group by ID
func UpdateGroup(c *gin.Context) {
	var group models.Group
	id := c.Param("id")

	if err := initializers.DBConn.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initializers.DBConn.Save(&group)
	c.JSON(http.StatusOK, group)
}

// DeleteGroup handles deleting a group by ID
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")

	if err := initializers.DBConn.Delete(&models.Group{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

// ListGroups retrieves all groups
func ListGroups(c *gin.Context) {
	var groups []models.Group

	initializers.DBConn.Preload("Members").Find(&groups)
	c.JSON(http.StatusOK, groups)
}
