package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getCatalog(c *gin.Context) {
	// path := strings.TrimPrefix(c.Param("path"), "/")
	// parts := strings.Split(path, "/")
	// catalogID := strings.TrimSuffix(parts[0], ".json")

	id := "1"
	catalog, err := s.store.GetCatalog("1", "_")

	if err != nil {
		slog.Error("failed to find catalog", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metas": catalog})
}
