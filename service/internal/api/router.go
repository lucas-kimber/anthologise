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

	// Public stremio routes
	r.GET("/:userID/manifest.json", server.getManifest)
	r.GET("/:userID/catalog/:type/:anthologyID", server.getCatalog)
	r.GET("/:userID/meta/:type/:anthologyID", server.getAnthology)

	r.GET("/health", server.getHealth)

	// Config site facing routes
	r.POST("/api/users", server.createUser)

	// Private api routes
	private := r.Group("/api/:userID")
	private.Use(server.authMiddleware())

	private.POST("/anthologies", server.createAnthology)
	private.PUT("/catalog/:anthologyID", server.addAnthologyToCatalog)
	private.DELETE("/catalog/:anthologyID", server.removeAnthologyFromCatalog)

	return r
}
