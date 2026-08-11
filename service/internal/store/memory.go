package store

import (
	"cmp"
	"log/slog"
	"slices"

	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type MemoryStore struct {
	catalogs    map[string]map[string]struct{}
	anthologies map[string]stremio.Anthology
}

var _ api.Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		catalogs:    make(map[string]map[string]struct{}),
		anthologies: make(map[string]stremio.Anthology),
	}
}

func (s *MemoryStore) ensureCatalog(token string) map[string]struct{} {
	c, ok := s.catalogs[token]

	if !ok {
		c = make(map[string]struct{})
		s.catalogs[token] = c

		slog.Info("created empty catalog")
	}

	return c
}

func (s *MemoryStore) GetAnthology(token string, anthologyID string) (stremio.Anthology, error) {

	a, ok := s.anthologies[anthologyID]

	if !ok {
		return stremio.Anthology{}, api.ErrAnthologyNotFound
	}

	return a, nil
}

func (s *MemoryStore) GetCatalog(token string) stremio.Catalog {
	c := s.ensureCatalog(token)

	metas := make([]stremio.AnthologyPreview, 0, len(c))

	for id := range c {
		anthology, ok := s.anthologies[id]
		if !ok {
			continue
		}

		metas = append(metas, anthology.AnthologyPreview)
	}

	slices.SortFunc(metas, func(i, j stremio.AnthologyPreview) int {
		return cmp.Compare(i.Name, j.Name)
	})

	return stremio.Catalog{
		Metas: metas,
	}
}

func (s *MemoryStore) AddAnthologyToCatalog(token string, anthologyID string) error {
	if _, ok := s.anthologies[anthologyID]; !ok {
		return api.ErrAnthologyNotFound
	}

	c := s.ensureCatalog(token)
	c[anthologyID] = struct{}{}

	return nil
}

func (s *MemoryStore) RemoveAnthologyFromCatalog(token string, anthologyID string) error {
	c := s.ensureCatalog(token)

	if _, ok := c[anthologyID]; !ok {
		return api.ErrAnthologyNotInCatalog
	}

	delete(c, anthologyID)

	return nil
}

func (s *MemoryStore) CreateAnthology(token string, anthology stremio.Anthology) error {

	if _, exists := s.anthologies[anthology.ID]; exists {
		return api.ErrAnthologyAlreadyExists
	}

	s.anthologies[anthology.ID] = anthology

	return nil
}

func (s *MemoryStore) UpdateAnthology(token string, anthology stremio.Anthology) error {
	if _, ok := s.anthologies[anthology.ID]; !ok {
		return api.ErrAnthologyNotFound
	}

	s.anthologies[anthology.ID] = anthology

	return nil
}
