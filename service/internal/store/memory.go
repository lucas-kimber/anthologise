package store

import (
	"log/slog"

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
	catalogs    map[string]stremio.Catalog
	anthologies map[anthologyKey]stremio.Anthology
}

var _ api.Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		catalogs:    make(map[string]stremio.Catalog),
		anthologies: make(map[anthologyKey]stremio.Anthology),
	}
}

func (s *MemoryStore) ensureCatalog(token string) stremio.Catalog {
	c, ok := s.catalogs[token]

	if !ok {

		c = stremio.Catalog{
			Metas: []stremio.AnthologyPreview{},
		}

		s.catalogs[token] = c

		slog.Info("created empty catalog")
	}

	return c
}

func (s *MemoryStore) AddAnthology(token string, anthology stremio.Anthology) error {

	c := s.ensureCatalog(token)
	c.Metas = append(c.Metas, anthology.AnthologyPreview)

	s.catalogs[token] = c

	ak := anthologyKey{token, anthology.ID}
	s.anthologies[ak] = anthology

	return nil
}

func (s *MemoryStore) GetCatalog(token string) stremio.Catalog {
	return s.ensureCatalog(token)
}

func (s *MemoryStore) GetAnthology(token string, anthologyID string) (stremio.Anthology, error) {

	k := anthologyKey{token, anthologyID}
	a, ok := s.anthologies[k]

	if !ok {
		return stremio.Anthology{}, api.ErrAnthologyNotFound
	}

	return a, nil
}
