package handler

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/service"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type PollHandler struct {
	PollService *service.PollService
}

func NewPollHandler(pollService *service.PollService) *PollHandler {
	return &PollHandler{PollService: pollService}
}

// CreatePoll godoc
// @Summary      Create a poll
// @Description  Create a new poll manually
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        poll body models.Poll true "Poll Details"
// @Success      200  {object}  models.Poll
// @Router       /polls [post]
func (h *PollHandler) CreatePoll(c *gin.Context) {
	var poll models.Poll
	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.PollService.Create(&poll)
	if err != nil {
		slog.Error("Failed to create poll", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Poll created successfully", "id", created.ID)

	c.JSON(http.StatusOK, created)
}

// GetPolls godoc
// @Summary      Get all polls
// @Description  Retrieve a list of all polls
// @Tags         polls
// @Produce      json
// @Success      200  {array}  models.Poll
// @Router       /polls [get]
func (h *PollHandler) GetPolls(c *gin.Context) {
	polls, err := h.PollService.GetPolls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, polls)
}

// GetPoll godoc
// @Summary      Get a poll
// @Description  Get a poll by its ID
// @Tags         polls
// @Produce      json
// @Param        id path string true "Poll ID"
// @Success      200  {object}  models.Poll
// @Router       /polls/{id} [get]
func (h *PollHandler) GetPoll(c *gin.Context) {
	id := c.Param("id")
	poll, err := h.PollService.GetPoll(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}
	c.JSON(http.StatusOK, poll)
}

// PutPoll godoc
// @Summary      Update a poll
// @Description  Update a poll details (used to resume or complete)
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        id path string true "Poll ID"
// @Param        poll body models.Poll true "Poll Details"
// @Success      200  {object}  models.Poll
// @Router       /polls/{id} [put]
func (h *PollHandler) PutPoll(c *gin.Context) {
	id := c.Param("id")
	var poll models.Poll
	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	poll.ID = id

	updated, err := h.PollService.UpdatePoll(&poll)
	if err != nil {
		slog.Error("Failed to update poll", "id", id, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Poll updated successfully", "id", id)

	c.JSON(http.StatusOK, updated)
}

// GetVotes godoc
// @Summary      Get all votes
// @Description  Get votes mapping the attendee to their choices
// @Tags         polls
// @Produce      json
// @Param        id path string true "Poll ID"
// @Success      200  {object}  map[string]map[string]int
// @Router       /polls/{id}/votes [get]
func (h *PollHandler) GetVotes(c *gin.Context) {
	id := c.Param("id")
	votes, err := h.PollService.GetVotes(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, votes)
}

// GetOutstanding godoc
// @Summary      Get outstanding voters
// @Description  Get attendees who have not yet voted
// @Tags         polls
// @Produce      json
// @Param        id path string true "Poll ID"
// @Success      200  {array}  string
// @Router       /polls/{id}/outstanding [get]
func (h *PollHandler) GetOutstanding(c *gin.Context) {
	id := c.Param("id")
	outstanding, err := h.PollService.GetOutstanding(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, outstanding)
}

// PutVote godoc
// @Summary      Submit a vote
// @Description  Submit a vote for an attendee
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        id path string true "Poll ID"
// @Param        attendee path string true "Attendee Name"
// @Param        choices body map[string]int true "Choices (-1, 0, or 1)"
// @Success      200  {object}  map[string]int
// @Router       /polls/{id}/votes/{attendee} [put]
func (h *PollHandler) PutVote(c *gin.Context) {
	id := c.Param("id")
	attendee := c.Param("attendee")

	var choices map[string]int
	if err := c.ShouldBindJSON(&choices); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	normalized, err := h.PollService.AddVote(id, attendee, choices)
	if err != nil {
		slog.Error("Failed to add vote", "id", id, "attendee", attendee, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Vote added successfully", "id", id, "attendee", attendee)

	c.JSON(http.StatusOK, normalized)
}

// GetResults godoc
// @Summary      Get poll results
// @Description  Get aggregated poll results
// @Tags         polls
// @Produce      json
// @Param        id path string true "Poll ID"
// @Success      200  {object}  map[string]int
// @Router       /polls/{id}/results [get]
func (h *PollHandler) GetResults(c *gin.Context) {
	id := c.Param("id")
	results, err := h.PollService.GetResults(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
