package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getCatalog(c *gin.Context) {

	token := c.Param("token")
	c.JSON(http.StatusOK, s.store.GetCatalog(token))
}
