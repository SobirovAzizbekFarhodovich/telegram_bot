package handler

import (
	"bot/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func buildResponse(c *gin.Context, status int, message string, reason string, data interface{}) {
	response := models.APIResponse{
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestURL: c.Request.URL.Path,
		Message:    message,
		Reason:     reason,
		Data:       data,
	}
	c.JSON(status, response)
}

// @Summary Create Password
// @Description Create a new password for the user
// @Tags Password
// @Accept json
// @Produce json
// @Param Password body models.Password true "Password"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /password [post]
func (h *HTTPHandler) CreatePassword(c *gin.Context) {
	var password models.Password
	if err := c.ShouldBindJSON(&password); err != nil {
		buildResponse(c, http.StatusBadRequest, "Invalid input", "Invalid JSON format", nil)
		return
	}
	if password.UserID == "" {
		buildResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found", nil)
		return
	}
	if err := h.service.PrService.CreatePassword(password.UserID, password); err != nil {
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
// @Param userID path string true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /password/{userID} [get]
func (h *HTTPHandler) GetAllPasswordsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		buildResponse(c, http.StatusBadRequest, "UserID required", "Missing userID parameter", nil)
		return
	}
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
// @Param userID query string true "User ID"
// @Param site query string true "Site name"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /password [get]
func (h *HTTPHandler) GetByName(c *gin.Context) {
	userID := c.Query("userID")
	site := c.Query("site")
	if userID == "" || site == "" {
		buildResponse(c, http.StatusBadRequest, "Invalid input", "User ID and Site name are required", nil)
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

