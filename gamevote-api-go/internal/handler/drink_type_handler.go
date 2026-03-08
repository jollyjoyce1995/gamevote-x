package handler

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/service"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type DrinkTypeHandler struct {
	DrinkTypeService *service.DrinkTypeService
}

func NewDrinkTypeHandler(service *service.DrinkTypeService) *DrinkTypeHandler {
	return &DrinkTypeHandler{DrinkTypeService: service}
}

// GetDrinkTypes godoc
// @Summary      Get preset drinks
// @Description  Get a list of all preset drinks
// @Tags         drinks
// @Produce      json
// @Success      200  {array}  models.DrinkType
// @Router       /drinks/presets [get]
func (h *DrinkTypeHandler) GetDrinkTypes(c *gin.Context) {
	types, err := h.DrinkTypeService.GetDrinkTypes()
	if err != nil {
		slog.Error("Failed to get drink types", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Retrieved drink types", "count", len(types))
	c.JSON(http.StatusOK, types)
}

// PostDrinkType godoc
// @Summary      Add custom drink preset
// @Description  Create a new custom drink preset saving it to the database
// @Tags         drinks
// @Accept       json
// @Produce      json
// @Param        drinkType body models.DrinkType true "Drink Type Details"
// @Success      200  {object}  models.DrinkType
// @Router       /drinks/presets [post]
func (h *DrinkTypeHandler) PostDrinkType(c *gin.Context) {
	var dt models.DrinkType
	if err := c.ShouldBindJSON(&dt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.DrinkTypeService.AddCustomDrinkType(&dt)
	if err != nil {
		slog.Error("Failed to add custom drink type", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Custom drink type added", "name", created.Name)

	c.JSON(http.StatusOK, created)
}
