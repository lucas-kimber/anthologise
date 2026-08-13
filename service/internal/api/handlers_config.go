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

	token := c.Param("userID")
	c.JSON(http.StatusOK, s.store.GetCatalog(token))
}

func (s *server) getAnthology(c *gin.Context) {

	anthologyID := c.Param("anthologyID")

	anthology, err := s.store.GetAnthology(anthologyID)
	if err != nil {
		slog.Error("failed to find anthology", "id", anthologyID)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.JSON(http.StatusOK, anthology)
}
