package handler

import (
	"bot/models"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) string {
	userID := c.GetString("telegram_user_id")
	return userID
}

func buildResponse(c *gin.Context, status int, message string, reason string, data interface{}) {
	response := gin.H{
		"timestamp":   time.Now().Format(time.RFC3339),
		"request_url": c.Request.URL.Path,
		"message":     message,
		"reason":      reason,
		"data":        data,
	}
	c.JSON(status, response)
}

// @Summary Create Password
// @Description Create a new password for the user
// @Tags Password
// @Accept json
// @Produce json
// @Param Password body models.Password true "Password"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /password [post]
func (h *HTTPHandler) CreatePassword(c *gin.Context) {
    userID := getUserID(c)
    if userID == "" {
        buildResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found", nil)
        return
    }
    var password models.Password
    if err := c.BindJSON(&password); err != nil {
        buildResponse(c, http.StatusBadRequest, "Invalid input", "Invalid JSON format", nil)
        return
    }
    password.UserID = userID
    if err := h.service.PrService.CreatePassword(userID, password);
        err != nil {
        buildResponse(c, http.StatusInternalServerError, "Failed to create password", "Database error", nil)
        return
    }
    buildResponse(c, http.StatusCreated, "Password created successfully", "", password)
}





// @Summary Get All Passwords
// @Description Get all passwords for a user by user ID
// @Tags Password
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /password/:userID [get]
func (h *HTTPHandler) GetAllPasswordsByUserID(c *gin.Context) {
	userID := getUserID(c)
	passwords, err := h.service.PrService.GetAllPasswordsByUserID(userID)
	if err != nil {
		buildResponse(c, http.StatusInternalServerError, "Failed to fetch passwords", "Database error", nil)
		return
	}
	if len(passwords) == 0 {
		buildResponse(c, http.StatusOK, "No passwords found", "Data not found", passwords)
		return
	}
	buildResponse(c, http.StatusOK, "Passwords retrieved successfully", "", passwords)
}

// @Summary Get Passwords By Site
// @Description Get passwords by site name for the user
// @Tags Password
// @Accept json
// @Produce json
// @Param site query string true "Site name"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /password [get]
func (h *HTTPHandler) GetByName(c *gin.Context) {
	userID := getUserID(c)
	site := c.Query("site")
	if site == "" {
		buildResponse(c, http.StatusBadRequest, "Invalid input", "Site name is required", nil)
		return
	}
	passwords, err := h.service.PrService.GetByName(userID, site)
	if err != nil {
		buildResponse(c, http.StatusInternalServerError, "Failed to fetch passwords", "Database error", nil)
		return
	}
	if len(passwords) == 0 {
		buildResponse(c, http.StatusOK, "No passwords found", "Data not found", passwords)
		return
	}
	buildResponse(c, http.StatusOK, "Passwords retrieved successfully", "", passwords)
}
