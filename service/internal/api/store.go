package api

import (
	"errors"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

var ErrNotFound = errors.New("store item not found")

// Store defines the methods that the API handlers expect to be available for retreiving resources from the database
type Store interface {
	GetCatalog(token string, catalogID string) (stremio.Catalog, error)
	GetAnthology(token string, catalogID string) (stremio.Anthology, error)
}
