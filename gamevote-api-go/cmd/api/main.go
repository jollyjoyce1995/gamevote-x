package main

import (
	"os"

	"gamevote-api-go/internal/handler"
	"gamevote-api-go/internal/logger"
	"gamevote-api-go/internal/service"
	"gamevote-api-go/internal/storage"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gamevote-api-go/docs" // Import swagger docs
)

// @title           GameVote API
// @version         1.0
// @description     API for the GameVote application to manage parties, polls, and votes.
// @host            localhost:8080
// @BasePath        /
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using default environment variables")
	}

	logger.Init()

	// Initialize Surreal DB
	slog.Info("Initializing Surreal database...")
	if err := storage.InitDB(); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer storage.CloseDB()

	// SSE Broker (must be created before services)
	broker := service.Broker

	// Repositories
	partyRepo := &storage.PartyRepository{}
	pollRepo := &storage.PollRepository{}
	beerRepo := &storage.BeerRepository{}
	voteRepo := &storage.VoteRepository{}
	userRepo := &storage.UserRepository{}
	drinkTypeRepo := &storage.DrinkTypeRepository{}
	gameRepo := &storage.GameRepository{}

	// Services
	pollService := service.NewPollService(pollRepo, voteRepo)
	partyService := service.NewPartyService(partyRepo, beerRepo, pollService, broker)
	userService := service.NewUserService(userRepo)
	drinkTypeService := service.NewDrinkTypeService(drinkTypeRepo)
	steamWorker := service.NewSteamWorker(gameRepo)

	// Seed presets & start background workers
	if err := drinkTypeService.SeedPresets(); err != nil {
		slog.Error("Failed to seed drink presets", "error", err)
	}
	steamWorker.Start()

	// Handlers
	partyHandler := handler.NewPartyHandler(partyService, broker)
	pollHandler := handler.NewPollHandler(pollService)
	userHandler := handler.NewUserHandler(userService)
	drinkTypeHandler := handler.NewDrinkTypeHandler(drinkTypeService)
	gameHandler := handler.NewGameHandler(steamWorker)

	router := gin.Default()

	// CORS for frontend dev server
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes
	router.POST("/users", userHandler.Login)
	router.GET("/users", userHandler.GetUsers)

	// Drink preset routes
	router.GET("/drinks/presets", drinkTypeHandler.GetDrinkTypes)
	router.POST("/drinks/presets", drinkTypeHandler.PostDrinkType)

	// Game search route
	router.GET("/games", gameHandler.SearchGames)

	// Party routes
	router.GET("/parties", partyHandler.GetParties)
	router.POST("/parties", partyHandler.CreateParty)
	router.GET("/parties/:code", partyHandler.GetParty)
	router.PATCH("/parties/:code", partyHandler.PatchParty)
	router.GET("/parties/:code/stream", partyHandler.StreamParty)
	router.POST("/parties/:code/options", partyHandler.PostOption)
	router.DELETE("/parties/:code/options/:gameName", partyHandler.DeleteOption)
	router.POST("/parties/:code/attendees", partyHandler.PostAttendee)
	router.DELETE("/parties/:code/attendees/:attendeeId", partyHandler.DeleteAttendee)
	router.POST("/parties/:code/beers", partyHandler.PostBeer)

	// Poll routes
	router.POST("/polls", pollHandler.CreatePoll)
	router.GET("/polls", pollHandler.GetPolls)
	router.GET("/polls/:id", pollHandler.GetPoll)
	router.PUT("/polls/:id", pollHandler.PutPoll)
	router.GET("/polls/:id/votes", pollHandler.GetVotes)
	router.GET("/polls/:id/outstanding", pollHandler.GetOutstanding)
	router.PUT("/polls/:id/votes/:attendee", pollHandler.PutVote)
	router.GET("/polls/:id/results", pollHandler.GetResults)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting server", "port", port)
	if err := router.Run("127.0.0.1:" + port); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
