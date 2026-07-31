package api

import "github.com/lucas-kimber/anthologise/service/internal/stremio"

// Store defines the methods that the API handlers expect to be available for retreiving resources from the database
type Store interface {
	GetManifest(token string) (stremio.Manifest, error)
	GetCatalog(token string, catalogID string) (stremio.Catalog, error)
	GetAnthology(token string, catalogID string) (stremio.Anthology, error)
}
