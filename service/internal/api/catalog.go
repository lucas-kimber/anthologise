package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *server) getCatalog(c *gin.Context) {

	token := c.Param("token")
	catalogID, _ := strings.CutSuffix(c.Param("id"), ".json")

	catalog, err := s.store.GetCatalog(token, catalogID)

	if err != nil {
		slog.Error("failed to find catalog", "id", catalogID)
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}

	c.JSON(http.StatusOK, catalog)
}
