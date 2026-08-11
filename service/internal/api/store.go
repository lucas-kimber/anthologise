package api

import (
	"errors"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

var ErrCatalogNotFound = errors.New("catalog item not found")
var ErrAnthologyNotFound = errors.New("anthology item not found")
var ErrAnthologyNotInCatalog = errors.New("anthology item not found")
var ErrAnthologyAlreadyExists = errors.New("anthology item already exists")

// Store defines the methods that the API handlers expect to be available for retreiving resources from the database
type Store interface {
	GetCatalog(token string) stremio.Catalog
	GetAnthology(token string, anthologyID string) (stremio.Anthology, error)
	CreateAnthology(token string, anthology stremio.Anthology) error
	UpdateAnthology(token string, anthology stremio.Anthology) error
	AddAnthologyToCatalog(token string, anthologyID string) error
}
