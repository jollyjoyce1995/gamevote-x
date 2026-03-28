package handler

import (
	"gamevote-api-go/internal/service"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	SteamWorker *service.SteamWorker
}

func NewGameHandler(steamWorker *service.SteamWorker) *GameHandler {
	return &GameHandler{SteamWorker: steamWorker}
}

// SearchGames godoc
// @Summary      Search games
// @Description  Search the cached Steam game list by name
// @Tags         games
// @Produce      json
// @Param        q query string true "Search Query"
// @Success      200  {array}  models.Game
// @ID           SearchGames
// @Router       /games [get]
func (h *GameHandler) SearchGames(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	games, err := h.SteamWorker.Search(q)
	if err != nil {
		slog.Error("Failed to search games", "query", q, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Game search returned results", "query", q, "count", len(games))

	c.JSON(http.StatusOK, games)
}
