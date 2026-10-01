package httpserver

import (
	"context"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

// Store defines the methods that the API handlers expect to be available for retreiving resources from the database
type Store interface {
	GetCatalog(ctx context.Context, userID string) (stremio.Catalog, error)
	GetAnthology(ctx context.Context, anthologyID string) (stremio.Anthology, error)

	Health(ctx context.Context) error

	CreateUser(ctx context.Context, userID string, editTokenHash [32]byte) error
	GetTokenHash(ctx context.Context, userID string) ([32]byte, error)

	CreateAnthology(ctx context.Context, anthology stremio.Anthology) error
	SetCatalog(ctx context.Context, userID string, anthologyIDs []string) error
}

type server struct {
	manifest stremio.Manifest
	store    Store
}

func newServer(manifest stremio.Manifest, store Store) *server {
	return &server{
		manifest: manifest,
		store:    store,
	}
}
