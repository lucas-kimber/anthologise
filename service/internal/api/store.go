package api

import (
	"errors"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

var ErrCatalogNotFound = errors.New("catalog item not found")
var ErrAnthologyNotFound = errors.New("anthology item not found")

// Store defines the methods that the API handlers expect to be available for retreiving resources from the database
type Store interface {
	GetCatalog(token string) stremio.Catalog
	GetAnthology(token string, catalogID string) (stremio.Anthology, error)
}
