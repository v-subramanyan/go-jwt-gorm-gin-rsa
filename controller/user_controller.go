package controller

import (
	"fmt"
	"jwt/initializers"
	"jwt/models"
	"jwt/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SigninData struct {
	Username string   `json:"username" bindings:"required"`
	Email    string   `json:"email" bindings:"required"`
	Password string   `json:"password" bindings:"required"`
	Roles    []string `json:"roles"`
	Groups   []string `json:"groups"`
}

type UpdateData struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Groups   []string `json:"groups"`
}

// CreateUser handles creating a new user
func CreateUser(c *gin.Context) {
	input := SigninData{}

	// Bind the input JSON to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before saving the user
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     input.Username,
		Email:    input.Email,
		Password: hashedPassword, // Store hashed password here
	}

	// Update roles by name
	if len(input.Roles) > 0 {
		var roles []models.Role
		for _, roleName := range input.Roles {
			var role models.Role
			if err := initializers.DBConn.Where("name = ?", roleName).First(&role).Error; err != nil {
				// Ensure roleName is treated as a string
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role name: " + string(roleName)})
				return
			}
			roles = append(roles, role)
		}
		user.Roles = roles
	}

	// Update groups by name
	if len(input.Groups) > 0 {
		var groups []models.Group
		for _, groupName := range input.Groups {
			var group models.Group
			if err := initializers.DBConn.Where("name = ?", groupName).First(&group).Error; err != nil {
				// Ensure groupName is treated as a string
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group name: " + string(groupName)})
				return
			}
			groups = append(groups, group)
		}
		user.Groups = groups
	}

	// Save the user in the database
	if err := initializers.DBConn.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate RSA keys with expiration
	privateKeyPEM, publicKeyPEM, expiresAt, err := utils.GenerateRSAKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate RSA keys"})
		return
	}

	// Save RSA keys in the database
	rsaKey := models.RSAKeyPair{
		PrivateKey: privateKeyPEM,
		PublicKey:  publicKeyPEM,
		UserID:     user.ID,
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	if err := initializers.DBConn.Create(&rsaKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save RSA keys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Name,
			"email":    user.Email,
		},
	})
}

// GetUser retrieves a single user by ID along with RSA key, groups, and roles
func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Preload related data
	if err := initializers.DBConn.Preload("Groups").Preload("Roles").Preload("RSAKey").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating a user by ID
func UpdateUser(c *gin.Context) {
	var input UpdateData
	var user models.User
	id := c.Param("id")

	// Fetch the existing user from the database
	if err := initializers.DBConn.Preload("Roles").Preload("Groups").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind the incoming JSON to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user fields
	user.Name = input.Username
	user.Email = input.Email

	// Update roles by name
	if len(input.Roles) > 0 {
		var roles []models.Role
		for _, roleName := range input.Roles {
			var role models.Role
			if err := initializers.DBConn.Where("name = ?", roleName).First(&role).Error; err != nil {
				// Ensure roleName is treated as a string
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role name: " + string(roleName)})
				return
			}
			roles = append(roles, role)
		}
		user.Roles = roles
	}

	// Update groups by name
	if len(input.Groups) > 0 {
		var groups []models.Group
		for _, groupName := range input.Groups {
			var group models.Group
			if err := initializers.DBConn.Where("name = ?", groupName).First(&group).Error; err != nil {
				// Ensure groupName is treated as a string
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group name: " + string(groupName)})
				return
			}
			groups = append(groups, group)
		}
		user.Groups = groups
	}

	// Save the updated user
	if err := initializers.DBConn.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles deleting a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := initializers.DBConn.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ListUsers retrieves all users along with their associated groups, roles, and RSA keys
func ListUsers(c *gin.Context) {
	var users []models.User

	// Preload related data and retrieve all users
	if err := initializers.DBConn.Preload("Groups").Preload("Roles").Preload("RSAKey").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ForceGen bool   `json:"forceTokenGen"`
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	// Find the user by email
	if err := initializers.DBConn.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the password
	if !utils.CheckPasswordHash(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if JWTToken already exists and is valid
	if user.JWTToken != "" && !input.ForceGen {
		// Fetch active RSA public key for the user
		var activeRSAKey models.RSAKeyPair
		if err := initializers.DBConn.Where("user_id = ? AND is_active = ?", user.ID, true).First(&activeRSAKey).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve RSA keys"})
			return
		}

		// Validate the existing token
		claims, err := ValidateJWT(user.JWTToken, string(activeRSAKey.PublicKey))
		if err == nil {
			// Token is valid, return the existing token
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful (existing token)",
				"token":   user.JWTToken,
				"claims":  claims, // Optionally return claims
			})
			return
		} else {
			// Log validation error for debugging
			fmt.Println("Token validation error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Existing token is invalid or expired"})
		}
	}

	// If JWT is invalid or expired, generate a new JWT token using the active RSA key
	var activeRSAKey models.RSAKeyPair
	if err := initializers.DBConn.Where("user_id = ? AND is_active = ?", user.ID, true).First(&activeRSAKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve RSA keys"})
		return
	}

	// Generate a new JWT token using the active RSA key
	token, err := GenerateJWT(user, activeRSAKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Save the new JWT token to the user
	user.JWTToken = token
	if err := initializers.DBConn.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the user with the new token"})
		return
	}

	// Return the JWT token as a response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
