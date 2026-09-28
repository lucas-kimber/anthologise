package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-kimber/anthologise/service/internal/store"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

const healthCheckTimeout = 2 * time.Second

func (s *server) getManifest(c *gin.Context) {
	c.JSON(http.StatusOK, s.manifest)
}

func (s *server) getCatalog(c *gin.Context) {

	userID := c.Param("userID")

	catalog, err := s.store.GetCatalog(c.Request.Context(), userID)

	if err != nil {
		slog.Error("failed to find catalog for user", "id", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}

	c.JSON(http.StatusOK, catalog)
}

func (s *server) getAnthology(c *gin.Context) {

	anthologyID := c.Param("anthologyID")

	anthology, err := s.store.GetAnthology(c.Request.Context(), anthologyID)
	if err != nil {
		slog.Error("failed to find anthology", "id", anthologyID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.JSON(http.StatusOK, anthology)
}

func (s *server) getHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), healthCheckTimeout)
	defer cancel()

	if err := s.store.Health(ctx); err != nil {
		slog.Error("health check failed", "error", err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	c.Status(http.StatusOK)
}

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

	id := stremio.AnthologyIDPrefix + uuid.NewString()
	anthology.ID = id

	if err := s.store.CreateAnthology(c.Request.Context(), anthology); err != nil {
		slog.Debug("could not create anthology in database", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create anthology"})
		return
	}

	c.JSON(http.StatusCreated, anthology)
}

func (s *server) setCatalog(c *gin.Context) {
	userID := c.Param("userID")

	var req struct {
		AnthologyIDs []string `json:"anthologyIDs"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid catalog"})

		return
	}

	if err := s.store.SetCatalog(c.Request.Context(), userID, req.AnthologyIDs); err != nil {

		if errors.Is(err, store.ErrAnthologyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "catalog contains a non-existant anthology id"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update catalog"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
