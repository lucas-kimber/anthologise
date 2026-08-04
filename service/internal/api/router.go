package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

// Creates a Gin router configured with Slog and the various
// anthologise API handlers.
func NewRouter(manifest stremio.Manifest, store Store, middleware ...gin.HandlerFunc) *gin.Engine {

	r := gin.New()
	r.Use(middleware...)
	r.Use(gin.Recovery())

	// Stremio requires CORS headers
	r.Use(cors.Default())

	server := newServer(manifest, store)

	r.GET("/:token/manifest.json", server.getManifest)
	r.GET("/:token/catalog/:type/:id", server.getCatalog)
	r.GET("/:token/meta/:type/:id", server.getMeta)

	return r
}
