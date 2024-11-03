package jwt

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shop/models"
)

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// Register new users
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
			Role     string `json:"role"` // Роль может быть указана явно, но необязательна
		}

		// Привязываем входные данные к структуре
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
			return
		}

		// Инициализируем нового пользователя
		var user models.User
		user.Email = input.Email

		// Хешируем пароль
		if err := user.HashPassword(input.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при шифровании пароля"})
			return
		}

		// Устанавливаем роль пользователя
		var role models.Role
		if input.Role == "" {
			// Если роль не указана, устанавливаем "user" по умолчанию
			if err := db.Where("name = ?", "user").First(&role).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Роль 'user' не найдена в системе"})
				return
			}
		} else {
			// Если роль указана, проверяем её наличие в базе данных
			if err := db.Where("name = ?", input.Role).First(&role).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Указанная роль не существует"})
				return
			}
		}

		// Присваиваем RoleID пользователю
		user.RoleID = role.ID

		// Сохраняем пользователя в базе данных
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
			return
		}

		// Ответ без пароля
		c.JSON(http.StatusCreated, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  role.Name,
		})
	}
}

// Login user and return JWT token
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := db.Preload("Role").Where("email = ?", loginData.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if err := user.CheckPassword(loginData.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		token, err := GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Возвращаем токен, user_id и роль
		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"user_id": user.ID,
			"role":    user.Role.Name, // Предполагаем, что роль доступна как `user.Role.Name`
		})
	}
}
