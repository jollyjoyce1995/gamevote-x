package handler

import (
	"gamevote-api-go/internal/service"
)

type PollHandler struct {
	PollService *service.PollService
}

func NewPollHandler(pollService *service.PollService) *PollHandler {
	return &PollHandler{PollService: pollService}
}

/*
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
*/
