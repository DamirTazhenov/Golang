package jwt

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"taskmanager/models"
)

// Register new users
func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// Декодируем запрос в структуру пользователя
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Валидация полей (например, с использованием пакета validator)
		if user.Mode != "individual" && user.Mode != "team" {
			http.Error(w, "Invalid mode: must be 'individual' or 'team'", http.StatusBadRequest)
			return
		}

		// Хэшируем пароль
		if err := user.HashPassword(user.Password); err != nil {
			http.Error(w, "Password encryption failed", http.StatusInternalServerError)
			return
		}

		// Сохраняем пользователя в базу данных
		if err := db.Create(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Возвращаем данные пользователя без пароля
		json.NewEncoder(w).Encode(user)
	}
}

// Login user and return JWT token
func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}

		if err := user.CheckPassword(loginData.Password); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		token, err := GenerateJWT(user.ID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
