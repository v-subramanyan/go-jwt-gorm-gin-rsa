package middleware

import (
	"jwt/controller"
	"jwt/initializers"
	"jwt/models"
	"net/http"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AdminRequired middleware to protect admin routes
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			log.Println("Unauthorized: Missing or invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Decode the token to extract claims without validating
		claims, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// We are not validating the token here; just decoding it
			return nil, nil
		})

		// Extract user ID from claims
		userIDFloat, ok := claims.Claims.(jwt.MapClaims)["user_id"].(float64)
		if !ok {
			log.Println("Unauthorized: Invalid user ID in claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// Fetch the user along with the RSA public key and roles
		var user models.User
		if err := initializers.DBConn.Preload("RSAKey").Preload("Roles").First(&user, userID).Error; err != nil {
			log.Println("Unauthorized: User not found:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if the user has a valid public key
		if user.RSAKey.PublicKey == "" {
			log.Println("Unauthorized: User does not have a valid public key")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate the token using the user's public key
		validatedClaims, err := controller.ValidateJWT(tokenString, string(user.RSAKey.PublicKey)) // Validate with actual public key
		if err != nil {
			log.Println("Unauthorized: Token validation with public key failed:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Here, you can use the validatedClaims if needed
		log.Println(validatedClaims.GetIssuer())
		// For instance, you could log or check additional claims.

		// Check if the user has the Admin role
		for _, role := range user.Roles {
			if role.Name == "admin" {
				c.Next() // User is authorized, proceed to the next handler
				return
			}
		}

		// If no Admin role is found, deny access
		log.Println("Forbidden: User does not have admin role")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}
