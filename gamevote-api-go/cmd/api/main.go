package main

import (
	"log"
	"os"

	"gamevote-api-go/internal/handler"
	"gamevote-api-go/internal/service"
	"gamevote-api-go/internal/storage"

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
		log.Println("No .env file found, using default environment variables")
	}
	// Initialize Surreal DB
	log.Println("Initializing Surreal database...")
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.CloseDB()
	defer storage.CloseDB()

	// Repositories
	partyRepo := &storage.PartyRepository{}
	pollRepo := &storage.PollRepository{}
	beerRepo := &storage.BeerRepository{}
	voteRepo := &storage.VoteRepository{}

	// Services
	pollService := service.NewPollService(pollRepo, voteRepo)
	partyService := service.NewPartyService(partyRepo, beerRepo, pollService)

	// Handlers
	partyHandler := handler.NewPartyHandler(partyService)
	pollHandler := handler.NewPollHandler(pollService)

	router := gin.Default()

	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	router.POST("/parties", partyHandler.CreateParty)
	router.GET("/parties/:code", partyHandler.GetParty)
	router.PATCH("/parties/:code", partyHandler.PatchParty)
	router.GET("/parties/:code/options", partyHandler.GetOptions)
	router.POST("/parties/:code/options", partyHandler.PostOption)
	router.DELETE("/parties/:code/options/:optionId", partyHandler.DeleteOption)
	router.GET("/parties/:code/attendees", partyHandler.GetAttendees)
	router.POST("/parties/:code/attendees", partyHandler.PostAttendee)
	router.DELETE("/parties/:code/attendees/:attendeeId", partyHandler.DeleteAttendee)
	router.POST("/parties/:code/beers", partyHandler.PostBeer)

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

	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
