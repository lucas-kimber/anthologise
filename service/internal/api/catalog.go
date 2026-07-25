package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	store "github.com/lucas-kimber/anthologise/service/internal/storage"
)

func getCatalog(c *gin.Context) {
	// path := strings.TrimPrefix(c.Param("path"), "/")
	// parts := strings.Split(path, "/")
	// catalogID := strings.TrimSuffix(parts[0], ".json")

	id := "1"
	catalog, ok := store.GetCatalogsByToken("1")

	if !ok {
		slog.Error("failed to find catalog", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metas": catalog})
}
