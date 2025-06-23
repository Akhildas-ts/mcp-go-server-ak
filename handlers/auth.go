package handlers

import (
	"mcp-go-server/models"
	"mcp-go-server/response"
	"mcp-go-server/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Login initiates GitHub OAuth login flow
func Login(c *gin.Context) {
	url, err := usecase.GetLoginURL()
	if err != nil {
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Failed to generate login URL", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	c.Redirect(http.StatusFound, url)
}

// Callback handles GitHub OAuth callback
func Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Authorization code is required", nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var callbackData models.CallbackData
	callbackData.Code = code

	// Validate the input
	if err := validator.New().Struct(callbackData); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Invalid callback data", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Process the callback
	session, err := usecase.HandleCallback(callbackData)
	if err != nil {
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Authentication failed", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Authentication successful", session, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID, exists := c.Get(models.UserIDKey)
	if !exists {
		errRes := response.ErrorClientResponse(http.StatusUnauthorized, "User not authenticated", nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	profile, err := usecase.GetUserProfile(userID.(string))
	if err != nil {
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Failed to retrieve profile", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Profile retrieved successfully", profile, nil)
	c.JSON(http.StatusOK, successRes)
}
