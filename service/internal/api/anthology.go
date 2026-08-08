package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func (s *server) getAnthology(c *gin.Context) {

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

func (s *server) addAnthology(c *gin.Context) {

	token := c.Param("token")

	var anthology stremio.Anthology

	if err := c.ShouldBindJSON(&anthology); err != nil {
		slog.Debug("invalid anthology request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anthology"})
		return
	}

	s.store.AddAnthology(token, anthology)

	c.Status(http.StatusCreated)
}
