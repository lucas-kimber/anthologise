package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getMeta(c *gin.Context) {
	id := "1"
	anthology, err := s.store.GetAnthology(id, "_")

	if err != nil {
		slog.Error("failed to find anthology", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": anthology})
}
