package handlers

import (
    "net/http"

    "github.com/DavutcanJ/customer-segmentation-api/internal/services"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler() *UserHandler {
    return &UserHandler{
        userService: services.NewUserService(),
    }
}

// @Summary Get user profile
// @Description Get current user profile information
// @Tags user
// @Produce json
// @Success 200 {object} models.UserProfile
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /api/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    // Get user ID from JWT token
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
        return
    }

    user, err := h.userService.GetUserByID(userID.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    userProfile := gin.H{
        "id":        user.ID.Hex(),
        "username":  user.Username,
        "email":     user.Email,
        "full_name": user.FullName,
        "role":      user.Role,
        "is_active": user.IsActive,
    }

    c.JSON(http.StatusOK, userProfile)
}