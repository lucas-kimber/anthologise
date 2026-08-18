package store

import (
	"cmp"
	"context"
	"log/slog"
	"slices"

	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type MemoryStore struct {
	users       map[string][32]byte
	catalogs    map[string]map[string]struct{}
	anthologies map[string]stremio.Anthology
}

var _ api.Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:       make(map[string][32]byte),
		catalogs:    make(map[string]map[string]struct{}),
		anthologies: make(map[string]stremio.Anthology),
	}
}

func (s *MemoryStore) ensureCatalog(userID string) map[string]struct{} {
	c, ok := s.catalogs[userID]

	if !ok {
		c = make(map[string]struct{})
		s.catalogs[userID] = c

		slog.Info("created empty catalog")
	}

	return c
}

func (s *MemoryStore) GetAnthology(ctx context.Context, anthologyID string) (stremio.Anthology, error) {

	a, ok := s.anthologies[anthologyID]

	if !ok {
		return stremio.Anthology{}, api.ErrAnthologyNotFound
	}

	return a, nil
}

func (s *MemoryStore) GetCatalog(ctx context.Context, userID string) (stremio.Catalog, error) {
	c := s.ensureCatalog(userID)

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
	}, nil
}

func (s *MemoryStore) AddAnthologyToCatalog(ctx context.Context, userID string, anthologyID string) error {
	if _, ok := s.anthologies[anthologyID]; !ok {
		return api.ErrAnthologyNotFound
	}

	c := s.ensureCatalog(userID)
	c[anthologyID] = struct{}{}

	return nil
}

func (s *MemoryStore) RemoveAnthologyFromCatalog(ctx context.Context, userID string, anthologyID string) error {
	c := s.ensureCatalog(userID)

	if _, ok := c[anthologyID]; !ok {
		return api.ErrAnthologyNotInCatalog
	}

	delete(c, anthologyID)

	return nil
}

func (s *MemoryStore) CreateAnthology(ctx context.Context, anthology stremio.Anthology) error {

	if _, exists := s.anthologies[anthology.ID]; exists {
		return api.ErrAnthologyAlreadyExists
	}

	s.anthologies[anthology.ID] = anthology

	return nil
}

func (s *MemoryStore) UpdateAnthology(ctx context.Context, userID string, anthology stremio.Anthology) error {
	if _, ok := s.anthologies[anthology.ID]; !ok {
		return api.ErrAnthologyNotFound
	}

	s.anthologies[anthology.ID] = anthology

	return nil
}

func (s *MemoryStore) CreateUser(ctx context.Context, userID string, editTokenHash [32]byte) error {

	if _, exists := s.users[userID]; exists {
		return api.ErrUserAlreadyExists
	}

	s.users[userID] = editTokenHash

	return nil
}

func (s *MemoryStore) GetTokenHash(ctx context.Context, userID string) ([32]byte, error) {

	tokenHash, ok := s.users[userID]

	if !ok {
		var b [32]byte
		return b, api.ErrUserDoesNotExist
	}

	return tokenHash, nil
}
