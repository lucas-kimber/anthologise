package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	store "github.com/lucas-kimber/anthologise/internal/storage"
)

func getMeta(c *gin.Context) {
	id := "1"
	anthology, ok := store.GetAnthologiesByToken(id)

	if !ok {
		slog.Error("failed to find anthology", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": anthology})
}
