package api

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

func createUserID() string {
	return uuid.NewString()
}

func creatEditToken() (string, [32]byte) {

	b := make([]byte, 32)
	rand.Read(b)

	editToken := base64.RawURLEncoding.EncodeToString(b)
	tokenHash := sha256.Sum256([]byte(editToken))

	return editToken, tokenHash
}

func (s *server) verifyUser(userID, editToken string) (bool, error) {

	targetHash, err := s.store.GetTokenHash(userID)

	if err != nil {
		return false, err
	}

	givenHash := sha256.Sum256([]byte(editToken))

	return subtle.ConstantTimeCompare(targetHash[:], givenHash[:]) == 1, nil
}

func (s *server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		editToken, err := c.Cookie("anthologise_edit_token")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		valid, err := s.verifyUser(userID, editToken)

		if !valid || err != nil {
			slog.Debug("user authentication failed", "userID", userID, "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
