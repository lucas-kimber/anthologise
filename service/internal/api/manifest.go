package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	store "github.com/lucas-kimber/anthologise/service/internal/storage"
)

func getManifest(c *gin.Context) {
	id := "1"
	manifest, ok := store.GetManifestByToken(id)

	if !ok {
		slog.Error("failed to find manifest", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "manifest not found"})
		return
	}

	c.JSON(http.StatusOK, manifest)
}
