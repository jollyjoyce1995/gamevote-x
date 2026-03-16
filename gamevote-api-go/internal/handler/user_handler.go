package handler

import (
	"gamevote-api-go/internal/service"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
}

// Login godoc
// @Summary      Login or Register a User
// @Description  Logs in a user by username. If they do not exist, they are created.
// @Tags         users
// @Accept       application/json
// @Produce      application/json
// @Param        req body UserLoginRequest true "Login Request"
// @Success      200 {object} service.UserDTO
// @Router       /users [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserService.LoginOrRegister(req.Username)
	if err != nil {
		slog.Error("Failed to login or register user", "username", req.Username, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("User logged in", "username", user.Username, "id", user.ID)

	c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Returns a list of all registered users.
// @Tags         users
// @Produce      application/json
// @Success      200  {array} service.UserDTO
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.UserService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
