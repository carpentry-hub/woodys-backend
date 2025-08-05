package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/carpentry-hub/woodys-backend/internal/config"
	"github.com/carpentry-hub/woodys-backend/internal/database"
	"github.com/carpentry-hub/woodys-backend/internal/handlers"
	"github.com/carpentry-hub/woodys-backend/internal/middleware"
	"github.com/carpentry-hub/woodys-backend/internal/repositories"
	"github.com/carpentry-hub/woodys-backend/internal/services"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Run database migrations
	// if err := db.AutoMigrate(); err != nil {
	// 	log.Fatalf("Failed to run database migrations: %v", err)
	// }

	// Initialize repositories
	repos := &repositories.Repositories{
		User:        repositories.NewUserRepository(db.GetDB()),
		Project:     repositories.NewProjectRepository(db.GetDB()),
		Comment:     repositories.NewCommentRepository(db.GetDB()),
		Rating:      repositories.NewRatingRepository(db.GetDB()),
		ProjectList: repositories.NewProjectListRepository(db.GetDB()),
	}

	// Initialize services
	srvs := &services.Services{
		User:        services.NewUserService(repos.User, repos.Project),
		Project:     services.NewProjectService(repos.Project, repos.User),
		Comment:     services.NewCommentService(repos.Comment, repos.User, repos.Project),
		Rating:      services.NewRatingService(repos.Rating, repos.Project, repos.User),
		ProjectList: services.NewProjectListService(repos.ProjectList, repos.Project, repos.User),
	}

	// Initialize handlers
	h := handlers.NewHandlers(srvs)

	// Setup router
	router := setupRouter(h)

	// Setup server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(h *handlers.Handlers) *mux.Router {
	router := mux.NewRouter()

	// Global middleware
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.Security)
	router.Use(middleware.CORS())

	// Rate limiting middleware (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	router.Use(rateLimiter.Middleware)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/", h.Home).Methods("GET")

	// User routes
	userRouter := api.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", h.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}", h.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", h.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id:[0-9]+}", h.DeleteUser).Methods("DELETE")
	userRouter.HandleFunc("/{id:[0-9]+}/projects", h.GetUserProjects).Methods("GET")
	userRouter.HandleFunc("/uid/{firebase_uid}", h.GetUserByUID).Methods("GET")

	// Project routes
	projectRouter := api.PathPrefix("/projects").Subrouter()
	projectRouter.HandleFunc("", h.CreateProject).Methods("POST")
	projectRouter.HandleFunc("/search", h.SearchProjects).Methods("GET")
	projectRouter.HandleFunc("/{id:[0-9]+}", h.GetProject).Methods("GET")
	projectRouter.HandleFunc("/{id:[0-9]+}", h.UpdateProject).Methods("PUT")
	projectRouter.HandleFunc("/{id:[0-9]+}", h.DeleteProject).Methods("DELETE")

	// Comment routes
	commentRouter := api.PathPrefix("/projects/{project_id:[0-9]+}/comments").Subrouter()
	commentRouter.HandleFunc("", h.GetProjectComments).Methods("GET")
	commentRouter.HandleFunc("", h.CreateComment).Methods("POST")

	commentManageRouter := api.PathPrefix("/comments").Subrouter()
	commentManageRouter.HandleFunc("/{id:[0-9]+}", h.DeleteComment).Methods("DELETE")
	commentManageRouter.HandleFunc("/{id:[0-9]+}/replies", h.GetCommentReplies).Methods("GET")
	commentManageRouter.HandleFunc("/{id:[0-9]+}/reply", h.CreateReply).Methods("POST")

	// Rating routes
	ratingRouter := api.PathPrefix("/projects/{project_id:[0-9]+}/ratings").Subrouter()
	ratingRouter.HandleFunc("", h.CreateRating).Methods("POST")
	ratingRouter.HandleFunc("", h.UpdateRating).Methods("PUT")
	ratingRouter.HandleFunc("", h.GetProjectRatings).Methods("GET")

	// Project list routes
	listRouter := api.PathPrefix("/project-lists").Subrouter()
	listRouter.HandleFunc("", h.CreateProjectList).Methods("POST")
	listRouter.HandleFunc("/{id:[0-9]+}", h.GetProjectList).Methods("GET")
	listRouter.HandleFunc("/{id:[0-9]+}", h.UpdateProjectList).Methods("PUT")
	listRouter.HandleFunc("/{id:[0-9]+}", h.DeleteProjectList).Methods("DELETE")
	listRouter.HandleFunc("/{id:[0-9]+}/projects", h.AddProjectToList).Methods("POST")
	listRouter.HandleFunc("/{list_id:[0-9]+}/projects/{project_id:[0-9]+}", h.RemoveProjectFromList).Methods("DELETE")

	userListRouter := api.PathPrefix("/users/{user_id:[0-9]+}/project-lists").Subrouter()
	userListRouter.HandleFunc("", h.GetUserProjectLists).Methods("GET")

	// 404 handler
	router.NotFoundHandler = http.HandlerFunc(middleware.NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(middleware.MethodNotAllowed)

	return router
}
