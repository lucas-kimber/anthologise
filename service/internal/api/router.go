package api

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// Creates a Gin router configured with Slog and the various
// anthologise API handlers.
func NewRouter() *gin.Engine {

	l := slog.Default()

	r := gin.New()
	r.Use(sloggin.New(l))
	r.Use(gin.Recovery())
	// Stremio requires CORS headers
	r.Use(cors.Default())

	r.GET("/manifest.json", getManifest)
	r.GET("/catalog/:type/*path", getCatalog)
	r.GET("/meta/:type/:id", getMeta)

	return r
}
