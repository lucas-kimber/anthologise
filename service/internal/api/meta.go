package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getMeta(c *gin.Context) {

	token := c.Param("token")
	id := c.Param("id")

	anthology, err := s.store.GetAnthology(token, id)
	if err != nil {
		slog.Error("failed to find anthology", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "anthology not found"})
		return
	}

	c.JSON(http.StatusOK, anthology)
}
