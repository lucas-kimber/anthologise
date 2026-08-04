package store

import (
	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type catalogKey struct {
	token     string
	catalogID string
}

type anthologyKey struct {
	token       string
	anthologyID string
}

type MemoryStore struct {
	catalogs    map[catalogKey]stremio.Catalog
	anthologies map[anthologyKey]stremio.Anthology
}

var _ api.Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		catalogs:    make(map[catalogKey]stremio.Catalog),
		anthologies: make(map[anthologyKey]stremio.Anthology),
	}
}

func (s *MemoryStore) AddCatalog(token string, catalogID string, catalog stremio.Catalog) {

	k := catalogKey{token, catalogID}
	s.catalogs[k] = catalog
}

func (s *MemoryStore) GetCatalog(token string, catalogID string) (stremio.Catalog, error) {

	k := catalogKey{token, catalogID}
	c, ok := s.catalogs[k]

	if !ok {
		return stremio.Catalog{}, api.ErrNotFound
	}

	return c, nil
}

func (s *MemoryStore) GetAnthology(token string, catalogID string) (stremio.Anthology, error) {

	k := anthologyKey{token, catalogID}
	a, ok := s.anthologies[k]

	if !ok {
		return stremio.Anthology{}, api.ErrNotFound
	}

	return a, nil
}
