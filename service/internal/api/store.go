package api

import (
	"errors"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

var ErrCatalogNotFound = errors.New("catalog item not found")
var ErrAnthologyNotFound = errors.New("anthology item not found")
var ErrAnthologyNotInCatalog = errors.New("anthology item not found")
var ErrAnthologyAlreadyExists = errors.New("anthology item already exists")
var ErrUserDoesNotExist = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user id already exists")

// Store defines the methods that the API handlers expect to be available for retreiving resources from the database
type Store interface {
	GetCatalog(userID string) stremio.Catalog
	GetAnthology(anthologyID string) (stremio.Anthology, error)

	CreateUser(userID string, editTokenHash [32]byte) error
	GetTokenHash(userID string) ([32]byte, error)

	CreateAnthology(userID string, anthology stremio.Anthology) error
	UpdateAnthology(userID string, anthology stremio.Anthology) error
	AddAnthologyToCatalog(userID string, anthologyID string) error
	RemoveAnthologyFromCatalog(userID string, anthologyID string) error
}
