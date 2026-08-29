package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func (s *server) createUser(c *gin.Context) {

	newID := createUserID()
	editToken, tokenHash := createEditToken()

	if err := s.store.CreateUser(c.Request.Context(), newID, tokenHash); err != nil {

		slog.Error("failed to create user", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	setEditTokenCookie(c, editToken)

	c.JSON(http.StatusCreated, gin.H{"userID": newID})
}

func (s *server) createAnthology(c *gin.Context) {

	var anthology stremio.Anthology

	if err := c.ShouldBindJSON(&anthology); err != nil {
		slog.Debug("invalid anthology request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anthology"})
		return
	}

	id := stremio.AnthologyIDPrefix + "_" + uuid.NewString()
	anthology.ID = id

	if err := s.store.CreateAnthology(c.Request.Context(), anthology); err != nil {
		slog.Debug("could not create anthology in database", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create anthology"})
		return
	}

	c.JSON(http.StatusCreated, anthology)
}

func (s *server) addAnthologyToCatalog(c *gin.Context) {

	userID := c.Param("userID")
	anthologyID := c.Param("anthologyID")

	if err := s.store.AddAnthologyToCatalog(c.Request.Context(), userID, anthologyID); err != nil {

		slog.Info("failed to find anthology", "id", anthologyID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *server) removeAnthologyFromCatalog(c *gin.Context) {
	userID := c.Param("userID")
	anthologyID := c.Param("anthologyID")

	if err := s.store.RemoveAnthologyFromCatalog(c.Request.Context(), userID, anthologyID); err != nil {

		slog.Info("failed to find anthology", "id", anthologyID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
