package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shop/models"
	"strconv"
)

var items = []models.Item{}

func createItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newItem models.Item

		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Get the user_id from the context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		newItem.UserID = userID.(uint)

		if err := db.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
			return
		}

		c.JSON(http.StatusCreated, newItem)
	}
}

func getItemByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var item models.Item
		if err := db.First(&item, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func getItems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.Item

		if err := db.Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
			return
		}

		c.JSON(http.StatusOK, items)
	}
}

func updateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		var item models.Item
		if err := db.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}

		if item.UserID != userID.(uint) && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this item"})
			return
		}

		var updatedItem models.Item
		if err := c.ShouldBindJSON(&updatedItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		item.Name = updatedItem.Name
		item.Price = updatedItem.Price
		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func deleteItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		var item models.Item
		if err := db.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}

		if item.UserID != userID.(uint) && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this item"})
			return
		}

		if err := db.Delete(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
	}
}
