package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func (s *server) createUser(c *gin.Context) {

	newID := createUserID()
	editToken, tokenHash := createEditToken()

	if err := s.store.CreateUser(newID, tokenHash); err != nil {

		slog.Error("failed to create user", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	setEditTokenCookie(c, editToken)

	c.JSON(http.StatusCreated, gin.H{"userID": newID})
}

func (s *server) verifyUser(userID, editToken string) (bool, error) {

	targetHash, err := s.store.GetTokenHash(userID)

	if err != nil {
		return false, err
	}

	givenHash := sha256.Sum256([]byte(editToken))

	return subtle.ConstantTimeCompare(targetHash[:], givenHash[:]) == 1, nil
}

func (s *server) addAnthology(c *gin.Context) {

	userID := c.Param("userID")

	var anthology stremio.Anthology

	if err := c.ShouldBindJSON(&anthology); err != nil {
		slog.Debug("invalid anthology request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anthology"})
		return
	}

	id := stremio.AnthologyIDPrefix + "_" + uuid.NewString()
	anthology.ID = id

	s.store.CreateAnthology(userID, anthology)
	s.store.AddAnthologyToCatalog(userID, anthology.ID)

	c.JSON(http.StatusCreated, anthology)
}

func (s *server) updateAnthology(c *gin.Context) {

	userID := c.Param("userID")

	var anthology stremio.Anthology

	if err := c.ShouldBindJSON(&anthology); err != nil {

		slog.Debug("invalid anthology request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anthology"})
		return
	}

	if err := s.store.UpdateAnthology(userID, anthology); err != nil {

		slog.Info("failed to find anthology", "id", anthology.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.Status(http.StatusOK)
}

func (s *server) addAnthologyToCatalog(c *gin.Context) {

	userID := c.Param("userID")
	anthologyID := c.Param("anthologyID")

	if err := s.store.AddAnthologyToCatalog(userID, anthologyID); err != nil {

		slog.Info("failed to find anthology", "id", anthologyID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}
}

func (s *server) removeAnthologyFromCatalog(c *gin.Context) {
	userID := c.Param("userID")
	anthologyID := c.Param("anthologyID")

	if err := s.store.RemoveAnthologyFromCatalog(userID, anthologyID); err != nil {

		slog.Info("failed to find anthology", "id", anthologyID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

}
