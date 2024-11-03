package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"shop/cmd/apiserver/jwt"
	"shop/internal/config"
	"shop/pkg/database"
)

var (
	DB         *gorm.DB
	configData *config.Config
)

func main() {
	database.InitDatabase()

	DB = database.GetDB()

	configData = config.NewConfig()

	router := gin.Default()

	router.POST("/register", func(c *gin.Context) { jwt.Register(DB)(c) })
	router.POST("/login", func(c *gin.Context) { jwt.Login(DB)(c) })

	api := router.Group("/api")
	api.Use(jwt.AuthMiddleware(DB))

	api.GET("/items", getItems(DB))
	router.GET("/api/items/:id", getItemByID(DB))
	api.POST("/items", createItem(DB))

	managerOrAdmin := api.Group("/")
	managerOrAdmin.Use(RequireRoles("manager", "admin"))
	managerOrAdmin.PUT("/items/:id", updateItem(DB))

	adminOnly := api.Group("/")
	adminOnly.Use(RequireRoles("admin"))
	adminOnly.DELETE("/items/:id", deleteItem(DB))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configData.BindAddr), handler))
}
