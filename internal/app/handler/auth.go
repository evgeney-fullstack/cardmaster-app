package handler

import (
	"github.com/gin-gonic/gin"
)

// signUp handles user registration
// Expected payload: {email, password, username}
// Returns: HTTP status 201 on success, error response on failure
func (h *Handler) signUp(c *gin.Context) {

}

// signIn handles user authentication
// Expected payload: {email, password}
// Returns: access token, refresh token, and user info on success
func (h *Handler) signIn(c *gin.Context) {

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
