package handler

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"github.com/gin-gonic/gin"
)

// signUp handles user registration
// Expected payload: {email, password, username}
// Returns: HTTP status 201 on success, error response on failure
func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	// Bind JSON request body to User model
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create user using authorization service
	err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return success response
	c.JSON(http.StatusOK, map[string]interface{}{
		"SignUp": "successful",
	})
}

func isValidEmail(email string) bool {
	// Проверка длины
	if len(email) > 254 {
		return false
	}

	// Парсинг email с помощью стандартной библиотеки
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Дополнительная проверка чтобы избежать неоднозначностей
	return addr.Address == email && !strings.Contains(email, " ")
}

// signIn handles user authentication
// Validates input credentials and returns JWT tokens on successful authentication
func (h *Handler) signIn(c *gin.Context) {
	var input models.SignInRequest

	// Bind and validate JSON input
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "error:invalid request format")
		return
	}

	// Validate that either email or username is provided
	if strings.TrimSpace(input.Email) == "" && strings.TrimSpace(input.Username) == "" {
		newErrorResponse(c, http.StatusBadRequest, "error:email or username is required")
		return
	}

	// Validate email format if provided
	if input.Email != "" && !isValidEmail(input.Email) {
		newErrorResponse(c, http.StatusBadRequest, "error:invalid email format")
		return
	}

	// Authenticate user and generate tokens
	authResult, err := h.services.Authorization.AuthenticateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "error:invalid credentials")
		return
	}

	// Return authentication response
	c.JSON(http.StatusOK, authResult)
}

// refreshHandler handles token refresh requests
// Expected header: Authorization: Bearer <refresh_token>
// Returns: new access token and refresh token
func (h *Handler) refreshTokens(c *gin.Context) {

}

// logout handles user logout (invalidates single refresh token)
// Expected header: Authorization: Bearer <refresh_token>
// Returns: HTTP status 200 on success
func (h *Handler) logout(c *gin.Context) {

}

// logoutAll handles user logout from all devices (invalidates all user's refresh tokens)
// Expected header: Authorization: Bearer <refresh_token>
// Returns: HTTP status 200 on success
func (h *Handler) logoutAll(c *gin.Context) {

}
