package handler

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PartyHandler struct {
	PartyService *service.PartyService
}

func NewPartyHandler(partyService *service.PartyService) *PartyHandler {
	return &PartyHandler{PartyService: partyService}
}

// CreateParty godoc
// @Summary      Create a new party
// @Description  Creates a new party with a generated 6-character code
// @Tags         parties
// @Accept       json
// @Produce      json
// @Param        party body models.Party true "Party details"
// @Success      200  {object}  service.PartyDTO
// @Router       /parties [post]
func (h *PartyHandler) CreateParty(c *gin.Context) {
	var party models.Party
	if err := c.ShouldBindJSON(&party); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.PartyService.CreateParty(&party)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dto, err := h.PartyService.ToDTO(created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetParty godoc
// @Summary      Get a party
// @Description  Get a party by its code
// @Tags         parties
// @Produce      json
// @Param        code path string true "Party Code"
// @Success      200  {object}  service.PartyDTO
// @Router       /parties/{code} [get]
func (h *PartyHandler) GetParty(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	party, err := h.PartyService.GetParty(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	dto, err := h.PartyService.ToDTO(party)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetOptions godoc
// @Summary      Get party options
// @Description  Get options for a specific party
// @Tags         parties
// @Produce      json
// @Param        code path string true "Party Code"
// @Success      200  {array}  string
// @Router       /parties/{code}/options [get]
func (h *PartyHandler) GetOptions(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	party, err := h.PartyService.GetParty(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	c.JSON(http.StatusOK, party.Options)
}

type StringValue struct {
	Value string `json:"value"`
}

// PostOption godoc
// @Summary      Add an option
// @Description  Add an option to a party
// @Tags         parties
// @Accept       json
// @Produce      json
// @Param        code path string true "Party Code"
// @Param        value body StringValue true "Option Value"
// @Success      200  {object}  StringValue
// @Router       /parties/{code}/options [post]
func (h *PartyHandler) PostOption(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	var val StringValue
	if err := c.ShouldBindJSON(&val); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.PartyService.AddOption(id, val.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, val)
}

// DeleteOption godoc
// @Summary      Delete an option
// @Description  Delete an option from a party by index
// @Tags         parties
// @Param        code path string true "Party Code"
// @Param        optionId path int true "Option Index"
// @Success      200
// @Router       /parties/{code}/options/{optionId} [delete]
func (h *PartyHandler) DeleteOption(c *gin.Context) {
	code := c.Param("code")
	optionIdStr := c.Param("optionId")
	optionId, err := strconv.Atoi(optionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option ID"})
		return
	}

	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	err = h.PartyService.DeleteOption(id, optionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetAttendees godoc
// @Summary      Get party attendees
// @Description  Get attendees for a specific party
// @Tags         parties
// @Produce      json
// @Param        code path string true "Party Code"
// @Success      200  {array}  string
// @Router       /parties/{code}/attendees [get]
func (h *PartyHandler) GetAttendees(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	party, err := h.PartyService.GetParty(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	c.JSON(http.StatusOK, party.Attendees)
}

// PostAttendee godoc
// @Summary      Add an attendee
// @Description  Add an attendee to a party
// @Tags         parties
// @Accept       json
// @Produce      json
// @Param        code path string true "Party Code"
// @Param        value body StringValue true "Attendee Name"
// @Success      200  {object}  StringValue
// @Router       /parties/{code}/attendees [post]
func (h *PartyHandler) PostAttendee(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	var val StringValue
	if err := c.ShouldBindJSON(&val); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.PartyService.AddAttendee(id, val.Value)
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
// @Router       /parties/{code}/attendees/{attendeeId} [delete]
func (h *PartyHandler) DeleteAttendee(c *gin.Context) {
	code := c.Param("code")
	attendeeIdStr := c.Param("attendeeId")
	attendeeId, err := strconv.Atoi(attendeeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	err = h.PartyService.DeleteAttendee(id, attendeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type PatchPartyRequest struct {
	Status string `json:"status"`
}

// PatchParty godoc
// @Summary      Patch a party
// @Description  Update a party status
// @Tags         parties
// @Accept       json
// @Produce      json
// @Param        code path string true "Party Code"
// @Param        patchReq body PatchPartyRequest true "Patch Request"
// @Success      200  {object}  service.PartyDTO
// @Router       /parties/{code} [patch]
func (h *PartyHandler) PatchParty(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	var req PatchPartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	party, err := h.PartyService.PatchParty(id, models.PartyStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto, err := h.PartyService.ToDTO(party)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

type BeerDTO struct {
	Attendee string `json:"attendee"`
}

// PostBeer godoc
// @Summary      Add a beer
// @Description  Add a beer for an attendee in a party
// @Tags         parties
// @Accept       json
// @Produce      json
// @Param        code path string true "Party Code"
// @Param        beer body BeerDTO true "Beer Details"
// @Success      200
// @Router       /parties/{code}/beers [post]
func (h *PartyHandler) PostBeer(c *gin.Context) {
	code := c.Param("code")
	id, err := h.PartyService.GetIdForCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	var beer BeerDTO
	if err := c.ShouldBindJSON(&beer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.PartyService.PostBeer(id, beer.Attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
