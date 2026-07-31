package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getManifest(c *gin.Context) {
	c.JSON(http.StatusOK, s.manifest)
}
