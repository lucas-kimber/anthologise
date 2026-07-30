package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getManifest(c *gin.Context) {
	id := "1"
	manifest, err := s.store.GetManifest(id)

	if err != nil {
		slog.Error("failed to find manifest", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "manifest not found"})
		return
	}

	c.JSON(http.StatusOK, manifest)
}
