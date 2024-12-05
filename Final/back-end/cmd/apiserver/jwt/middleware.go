package jwt

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shop/models"
	"strings"
)

// AuthMiddleware middleware for JWT authentication
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется заголовок Authorization"})
			c.Abort()
			return
		}

		token := strings.Split(authHeader, " ")
		if len(token) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильный формат заголовка Authorization"})
			c.Abort()
			return
		}

		claims, err := ValidateJWT(token[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.Preload("Role").First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", user.Role.Name)

		c.Next()
	}
}
