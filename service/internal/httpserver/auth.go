package httpserver

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const editTokenCookieName = "__Host-anthologise_edit_token"
const cookieMaxAge = 400 * 24 * 60 * 60

func createUserID() string {
	return uuid.NewString()
}

func createEditToken() (string, [32]byte) {

	b := make([]byte, 32)
	rand.Read(b)

	editToken := base64.RawURLEncoding.EncodeToString(b)
	tokenHash := sha256.Sum256([]byte(editToken))

	return editToken, tokenHash
}

func setEditTokenCookie(c *gin.Context, editToken string) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		editTokenCookieName,
		editToken,
		cookieMaxAge,
		"/",
		"",
		true,
		true,
	)
}

func (s *server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		editToken, err := c.Cookie(editTokenCookieName)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		targetHash, err := s.store.GetTokenHash(c.Request.Context(), userID)

		if err != nil {
			slog.Debug("user authentication failed, could not retreive user", "userID", userID, "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		givenHash := sha256.Sum256([]byte(editToken))
		valid := subtle.ConstantTimeCompare(targetHash[:], givenHash[:]) == 1

		if !valid {
			slog.Debug("user authentication failed", "userID", userID, "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Refresh the editToken cookie on successful auth
		setEditTokenCookie(c, editToken)
		c.Next()
	}
}
