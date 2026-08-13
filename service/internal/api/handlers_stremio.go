package api

import (
	"crypto/rand"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func (s *server) addAnthology(c *gin.Context) {

	userID := c.Param("userID")

	var anthology stremio.Anthology

	if err := c.ShouldBindJSON(&anthology); err != nil {
		slog.Debug("invalid anthology request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anthology"})
		return
	}

	id := stremio.AnthologyIDPrefix + rand.Text()
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
