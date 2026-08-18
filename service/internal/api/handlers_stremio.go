package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
