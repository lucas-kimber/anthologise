package api

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
	sloggin "github.com/samber/slog-gin"
)

// Creates a Gin router configured with Slog and the various
// anthologise API handlers.
func NewRouter(manifest stremio.Manifest, store Store) *gin.Engine {

	l := slog.Default()

	r := gin.New()
	r.Use(sloggin.New(l))
	r.Use(gin.Recovery())
	// Stremio requires CORS headers
	r.Use(cors.Default())

	server := newServer(manifest, store)

	r.GET("/manifest.json", server.getManifest)
	r.GET("/catalog/:type/*path", server.getCatalog)
	r.GET("/meta/:type/:id", server.getMeta)

	return r
}
