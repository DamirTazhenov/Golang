package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"taskmanager/cmd/apiserver/jwt"
	"taskmanager/internal/config"
	"taskmanager/pkg/database"
	"taskmanager/pkg/logging"
)

var (
	DB         *gorm.DB
	configData *config.Config
	logger     *logging.Logger
)

func main() {
	logger = logging.GetLogger()
	logger.Info("Starting server...")

	logger.Info("Initializing database...")
	database.InitDatabase()

	DB = database.GetDB()

	logger.Info("Initializing configuration...")
	configData = config.NewConfig()

	router := mux.NewRouter()

	CSRF := csrf.Protect(
		[]byte(configData.CSRFSecret),
		csrf.Secure(false),
		csrf.Path("/submit"),
	)

	router.Use(RequestLoggingMiddleware, MetricsMiddleware, ExpvarMiddleware)

	router.Handle("/metrics", promhttp.Handler())

	// Auth routes
	router.HandleFunc("/register", jwt.Register(DB)).Methods("POST")
	router.HandleFunc("/login", jwt.Login(DB)).Methods("POST")

	// Protected routes with prefix /api and jwt token
	api := router.PathPrefix("/api").Subrouter()
	api.Use(jwt.AuthMiddleware)

	// Индивидуальные задачи
	api.HandleFunc("/tasks", ModeMiddleware(GetAllTasks, nil)).Methods("GET")
	api.HandleFunc("/tasks", ModeMiddleware(CreateTask, nil)).Methods("POST")
	api.HandleFunc("/tasks/{id}", ModeMiddleware(GetTask, nil)).Methods("GET")
	api.HandleFunc("/tasks/{id}", ModeMiddleware(UpdateTask, nil)).Methods("PUT")
	api.HandleFunc("/tasks/{id}", ModeMiddleware(DeleteTask, nil)).Methods("DELETE")

	// Командные задачи
	api.Handle("/teams/{team_id}/tasks", TeamRoleMiddleware("manager", "employee", "client")(http.HandlerFunc(GetTeamTasks))).Methods("GET")
	api.Handle("/teams/{team_id}/members", TeamRoleMiddleware("manager")(http.HandlerFunc(AddTeamMember))).Methods("POST")

	// Team create
	api.HandleFunc("/teams", CreateTeam).Methods("POST")

	// CSRF-protected routes (apply middleware selectively)
	protected := router.PathPrefix("/form").Subrouter()
	protected.HandleFunc("/submit", SubmitHandler).Methods("POST")
	protected.Use(CSRF)

	// Example route for showing a form with the CSRF token
	router.HandleFunc("/form", FormHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Security Headers Middleware
	secureRouter := SecurityHeadersMiddleware(router)

	// Apply middleware (CORS, Security Headers, and CSRF)
	handler := c.Handler(secureRouter)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configData.BindAddr), handler))
}
