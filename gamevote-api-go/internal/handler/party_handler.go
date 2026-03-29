package handler

import (
	"fmt"
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/service"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type PartyHandler struct {
	PartyService *service.PartyService
	Broker       *service.SSEBroker
}

func NewPartyHandler(partyService *service.PartyService, broker *service.SSEBroker) *PartyHandler {
	return &PartyHandler{PartyService: partyService, Broker: broker}
}

// GetParties godoc
// @Summary      Get all parties
// @Description  Get all parties ordered by ID
// @Tags         parties
// @Produce application/json
// @Success      200  {array}  service.PartyDTO
// @ID           GetParties
// @Router       /parties [get]
func (h *PartyHandler) GetParties(c *gin.Context) {
	dto, err := h.PartyService.GetParties()
	if err != nil {
		slog.Error("Failed to get parties", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Successfully retrieved all parties", "count", len(dto))
	c.JSON(http.StatusOK, dto)
}

// StreamParty godoc
// @Summary	SSE stream for a party
// @Description	Opens a Server-Sent Events stream for real-time party updates
// @Tags parties
// @Produce	text/event-stream
// @Param code path string true "Party Code"
// @Param username query string true "Username of the connected client"
// @Success	200 {object} string
// @ID           StreamParty
// @Router       /parties/{code}/stream [get]
func (h *PartyHandler) StreamParty(c *gin.Context) {
	code := c.Param("code")
	username := c.Query("username")
	if username == "" {
		username = fmt.Sprintf("user-%d", time.Now().UnixNano())
	}
	clientID := username

	slog.Info("Opening SSE stream", "code", code, "username", username)
	client := h.Broker.Register(code, clientID)
	defer h.Broker.Unregister(code, clientID)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Send initial online users list
	onlineUsers := h.Broker.OnlineUsers(code)
	c.SSEvent("online_users", onlineUsers)
	c.Writer.Flush()

	// Send initial party state
	party, err := h.PartyService.GetPartyByCode(code)
	if err == nil {
		c.SSEvent("party_updated", party)
		c.Writer.Flush()
	}

	// Send initial outstanding voters if applicable
	domainParty, err := h.PartyService.PartyRepo.FindByCode(code)
	if err == nil && domainParty.Status == models.PartyStatusVoting {
		outstanding := h.PartyService.GetOutstandingVoters(domainParty)
		c.SSEvent("outstanding_voters_updated", outstanding)
		c.Writer.Flush()
	}

	slog.Info("Client connected to SSE", "code", code, "username", username)

	ctx := c.Request.Context()
	for {
		select {
		case msg, ok := <-client.Channel:
			if !ok {
				return
			}
			_, err := fmt.Fprint(c.Writer, msg)
			if err != nil {
				return
			}
			c.Writer.Flush()
		case <-ctx.Done():
			return
		}
	}
}

// CreateParty godoc
// @Summary      Create a new party
// @Description  Creates a new party with a generated 6-character code
// @Tags         parties
// @Accept       application/json
// @Produce      application/json
// @Param        party body service.PartyDTO true "Party details"
// @Success      200  {object}  service.PartyDTO
// @ID           CreateParty
// @Router       /parties [post]
func (h *PartyHandler) CreateParty(c *gin.Context) {
	var party service.PartyDTO
	if err := c.ShouldBindJSON(&party); err != nil {
		slog.Warn("Failed to bind JSON for party creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Creating new party", "creator", party.Attendees)
	partyDomain, err := h.PartyService.ToDomain(&party)
	if err != nil {
		slog.Error("Failed to convert party to domain", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.PartyService.CreateParty(partyDomain)
	if err != nil {
		slog.Error("Failed to create party", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Party created successfully", "code", created.Code)

	dto, err := h.PartyService.ToDTO(created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetParty godoc
// @Summary Get a party
// @Description Get a party by its code
// @Tags parties
// @Produce application/json
// @Param code path string true "Party Code"
// @Success 200 {object} service.PartyDTO
// @ID           GetParty
// @Router /parties/{code} [get]
func (h *PartyHandler) GetParty(c *gin.Context) {
	code := c.Param("code")
	dto, err := h.PartyService.GetPartyByCode(code)
	if err != nil {
		slog.Warn("Party not found", "code", code)
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	c.JSON(http.StatusOK, dto)
}

type StringValue struct {
	Value string `json:"value"`
}

// PostOption godoc
// @Summary      Add an option
// @Description  Add an option to a party
// @Tags         parties
// @Accept       application/json
// @Produce      application/json
// @Param        code path string true "Party Code"
// @Param        option body models.PartyOption true "Option Details"
// @Success      200  {object}  models.PartyOption
// @ID           PostOption
// @Router       /parties/{code}/options [post]
func (h *PartyHandler) PostOption(c *gin.Context) {
	code := c.Param("code")

	var opt models.PartyOption
	if err := c.ShouldBindJSON(&opt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Adding option", "code", code, "option", opt)
	err := h.PartyService.AddOption(code, opt)
	if err != nil {
		slog.Error("Failed to add option", "code", code, "option", opt.Name, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Option added", "code", code, "option", opt.Name)

	c.JSON(http.StatusOK, opt)
}

// DeleteOption godoc
// @Summary Delete an option
// @Description Delete an option from a party by its name
// @Tags parties
// @Param code path string true "Party Code"
// @Param gameName path string true "Game Name"
// @Success 200
// @ID           DeleteOption
// @Router /parties/{code}/options/{gameName} [delete]
func (h *PartyHandler) DeleteOption(c *gin.Context) {
	code := c.Param("code")
	gameName := c.Param("gameName")

	err := h.PartyService.DeleteOption(code, gameName)
	if err != nil {
		slog.Error("Failed to delete option", "code", code, "option", gameName, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Option deleted", "code", code, "option", gameName)

	c.Status(http.StatusOK)
}

// PostAttendee godoc
// @Summary      Add an attendee
// @Description  Add an attendee to a party
// @Tags         parties
// @Accept       application/json
// @Produce      application/json
// @Param        code path string true "Party Code"
// @Param        value body StringValue true "Attendee Name"
// @Success      200 {object} StringValue
// @ID           PostAttendee
// @Router       /parties/{code}/attendees [post]
func (h *PartyHandler) PostAttendee(c *gin.Context) {
	code := c.Param("code")
	var val StringValue
	if err := c.ShouldBindJSON(&val); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.PartyService.AddAttendee(code, val.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, val)
}

// DeleteAttendee godoc
// @Summary      Delete an attendee
// @Description  Delete an attendee from a party by index
// @Tags         parties
// @Param        code path string true "Party Code"
// @Param        attendeeId path int true "Attendee Index"
// @Success      200
// @ID           DeleteAttendee
// @Router       /parties/{code}/attendees/{attendeeId} [delete]
func (h *PartyHandler) DeleteAttendee(c *gin.Context) {
	code := c.Param("code")
	attendeeIdStr := c.Param("attendeeId")
	attendeeId, err := strconv.Atoi(attendeeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	err = h.PartyService.DeleteAttendee(code, attendeeId)
	if err != nil {
		slog.Error("Failed to delete attendee", "code", code, "attendeeId", attendeeId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Attendee deleted", "code", code, "attendeeId", attendeeId)

	c.Status(http.StatusOK)
}

type PatchPartyRequest struct {
	Status string `json:"status"`
}

// PatchParty godoc
// @Summary      Patch a party
// @Description  Update a party status
// @Tags         parties
// @Accept       application/json
// @Produce      application/json
// @Param        code path string true "Party Code"
// @Param        patchReq body PatchPartyRequest true "Patch Request"
// @Success      200 {object} service.PartyDTO
// @ID           PatchParty
// @Router       /parties/{code} [patch]
func (h *PartyHandler) PatchParty(c *gin.Context) {
	code := c.Param("code")

	var req PatchPartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to bind JSON for party patch", "code", code, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Patching party", "code", code, "status", req.Status)

	partyFromDb, err := h.PartyService.PatchParty(code, models.PartyStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto, err := h.PartyService.ToDTO(partyFromDb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// PostVote godoc
// @Summary      Submit a vote
// @Description  Submit a vote for an attendee
// @ID 			 PostVote
// @Tags         parties
// @Accept       application/json
// @Produce      application/json
// @Param        code path string true "Poll ID"
// @Param        attendee path string true "Attendee Name"
// @Param        choices body map[string]int true "Choices (-1, 0, or 1)"
// @Success      204
// @Router       /parties/{code}/votes/{attendee} [POST]
func (h *PartyHandler) PostVote(c *gin.Context) {
	code := c.Param("code")
	attendee := c.Param("attendee")

	var choices map[string]int
	if err := c.ShouldBindJSON(&choices); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.PartyService.AddPoll(code, attendee, choices)
	if err != nil {
		slog.Error("Failed to add vote", "code", code, "attendee", attendee, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Vote added successfully", "code", code, "attendee", attendee)

	c.Status(http.StatusNoContent)
}

type BeerDTO struct {
	Attendee string `json:"attendee"`
}

// PostBeer godoc
// @Summary      Add a beer
// @Description  Add a beer for an attendee in a party
// @Tags         parties
// @Accept       application/json
// @Produce      application/json
// @Param        code path string true "Party Code"
// @Param        beer body BeerDTO true "Beer Details"
// @Success      200
// @ID           PostBeer
// @Router       /parties/{code}/beers [post]
func (h *PartyHandler) PostBeer(c *gin.Context) {
	code := c.Param("code")

	var beer BeerDTO
	if err := c.ShouldBindJSON(&beer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.PartyService.PostBeer(code, beer.Attendee)
	if err != nil {
		slog.Error("Failed to add beer", "code", code, "attendee", beer.Attendee, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Beer added", "code", code, "attendee", beer.Attendee)

	c.Status(http.StatusOK)
}

// PostNextRound godoc
// @Summary      Raise the round
// @Description  Start the next round for a party
// @Tags         parties
// @Produce      application/json
// @Param        code path string true "Party Code"
// @Success      200 {object} service.PartyDTO
// @ID           PostNextRound
// @Router       /parties/{code}/next-round [post]
func (h *PartyHandler) PostNextRound(c *gin.Context) {
	code := c.Param("code")

	slog.Info("Starting next round", "code", code)

	partyFromDb, err := h.PartyService.NextRound(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto, err := h.PartyService.ToDTO(partyFromDb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}
