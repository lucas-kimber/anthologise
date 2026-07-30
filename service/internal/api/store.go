package api

import "github.com/lucas-kimber/anthologise/service/internal/stremio"

type Store interface {
	GetManifest(token string) (stremio.Manifest, error)
	GetCatalog(token string, catalogID string) (stremio.Catalog, error)
	GetAnthology(token string, catalogID string) (stremio.Anthology, error)
}
