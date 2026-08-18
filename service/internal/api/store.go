package api

import (
	"context"
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
	GetCatalog(ctx context.Context, userID string) (stremio.Catalog, error)
	GetAnthology(ctx context.Context, anthologyID string) (stremio.Anthology, error)

	CreateUser(ctx context.Context, userID string, editTokenHash [32]byte) error
	GetTokenHash(ctx context.Context, userID string) ([32]byte, error)

	CreateAnthology(ctx context.Context, anthology stremio.Anthology) error
	AddAnthologyToCatalog(ctx context.Context, userID string, anthologyID string) error
	RemoveAnthologyFromCatalog(ctx context.Context, userID string, anthologyID string) error
}
