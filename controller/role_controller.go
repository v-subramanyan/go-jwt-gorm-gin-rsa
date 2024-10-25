package controller

import (
	"jwt/initializers"
	"jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateRole handles creating a new role
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DBConn.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// GetRole retrieves a single role by ID
func GetRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := initializers.DBConn.Preload("Users").Preload("Groups").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

// UpdateRole handles updating a role by ID
func UpdateRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := initializers.DBConn.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initializers.DBConn.Save(&role)
	c.JSON(http.StatusOK, role)
}

// DeleteRole handles deleting a role by ID
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	if err := initializers.DBConn.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}

// ListRoles retrieves all roles
func ListRoles(c *gin.Context) {
	var roles []models.Role

	initializers.DBConn.Preload("Users").Preload("Groups").Find(&roles)
	c.JSON(http.StatusOK, roles)
}
