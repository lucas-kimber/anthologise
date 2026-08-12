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

	r.GET("/:userID/manifest.json", server.getManifest)
	r.GET("/:userID/catalog/:type/:anthologyID", server.getCatalog)
	r.GET("/:userID/meta/:type/:anthologyID", server.getAnthology)

	r.POST("/api/:userID/anthologies", server.addAnthology)
	r.PUT("/api/:userID/anthologies", server.updateAnthology)

	return r
}
