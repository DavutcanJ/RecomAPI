package handlers

import (
    "net/http"

    "github.com/DavutcanJ/customer-segmentation-api/internal/models"
    "github.com/DavutcanJ/customer-segmentation-api/internal/services"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    userService *services.UserService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        userService: services.NewUserService(),
    }
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.UserProfile
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
        return
    }

    user, err := h.userService.RegisterUser(req)
    if err != nil {
        statusCode := http.StatusBadRequest
        if err.Error() == "username already exists" || err.Error() == "email already exists" {
            statusCode = http.StatusConflict
        }
        c.JSON(statusCode, models.ErrorResponse{Error: err.Error()})
        return
    }

    userProfile := models.UserProfile{
        ID:       user.ID.Hex(),
        Username: user.Username,
        Email:    user.Email,
        FullName: user.FullName,
        Role:     user.Role,
        IsActive: user.IsActive,
    }

    c.JSON(http.StatusCreated, userProfile)
}

// @Summary User login
// @Description Login with username and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials"
// @Success 200 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
        return
    }

    loginResponse, err := h.userService.LoginUser(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, loginResponse)
}