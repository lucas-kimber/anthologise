package api_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/store"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func testAnthology(id string, name string) stremio.Anthology {
	return stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          id,
			Type:        "series",
			Name:        name,
			Poster:      "Test PosterURL",
			Description: "Test Description",
			Genres:      []string{"Test"},
		},
		Videos: []stremio.Video{
			{
				ID:       "test_video",
				Title:    "Test Video",
				Season:   1,
				Episode:  1,
				Released: "Test Released",
				Overview: "Test Overview",
			},
		},
	}
}

func seedCatalog(
	t *testing.T,
	s *store.MemoryStore,
	anthologies ...stremio.Anthology,
) {
	t.Helper()

	for _, anthology := range anthologies {
		if err := s.CreateAnthology(context.Background(), anthology); err != nil {
			t.Fatalf("failed to create anthology: %v", err)
		}

		if err := s.AddAnthologyToCatalog(
			context.Background(),
			testUserID,
			anthology.ID,
		); err != nil {
			t.Fatalf("failed to add anthology to catalog: %v", err)
		}
	}
}

func TestGetAnthology(t *testing.T) {
	want := testAnthology(
		"testid",
		"Test Anthology",
	)

	s := store.NewMemoryStore()

	if err := s.CreateAnthology(context.Background(), want); err != nil {
		t.Fatalf("failed to create anthology: %v", err)
	}

	router := newTestRouter(s)

	got := getJSON[stremio.Anthology](
		t,
		router,
		"/testUserID/meta/series/testid",
	)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect anthology:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}

func TestGetCatalog(t *testing.T) {
	first := testAnthology(
		"testid1",
		"Test Anthology 1",
	)
	second := testAnthology(
		"testid2",
		"Test Anthology 2",
	)

	s := store.NewMemoryStore()

	seedCatalog(
		t,
		s,
		first,
		second,
	)

	router := newTestRouter(s)

	got := getJSON[stremio.Catalog](
		t,
		router,
		"/testUserID/catalog/series/anthologise",
	)

	want := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			first.AnthologyPreview,
			second.AnthologyPreview,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect catalog:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}

func TestAddAnthology(t *testing.T) {
	want := testAnthology(
		"",
		"Test Anthology",
	)

	s := store.NewMemoryStore()
	seedUser(t, s)

	router := newTestRouter(s)

	res := authenticatedJSONRequest(
		t,
		router,
		http.MethodPost,
		"/api/testUserID/anthologies",
		want,
	)

	requireStatus(t, res, http.StatusCreated)

	created := decodeJSON[stremio.Anthology](t, res)

	if !strings.HasPrefix(created.ID, stremio.AnthologyIDPrefix) {
		t.Errorf("generated anthology ID has incorrect prefix: %q", created.ID)
	}

	want.ID = created.ID

	if !reflect.DeepEqual(created, want) {
		t.Errorf("created incorrect anthology:\nwant: %+v\ngot:  %+v", want, created)
	}

	got := getJSON[stremio.Anthology](
		t,
		router,
		"/testUserID/meta/series/"+created.ID,
	)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("stored incorrect anthology:\nwant: %+v\ngot:  %+v", want, got)
	}

	catalog := getJSON[stremio.Catalog](
		t,
		router,
		"/testUserID/catalog/series/anthologise",
	)

	wantCatalog := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			want.AnthologyPreview,
		},
	}

	if !reflect.DeepEqual(catalog, wantCatalog) {
		t.Errorf("received incorrect catalog:\nwant: %+v\ngot:  %+v", wantCatalog, catalog)
	}
}

func TestAddAnthologyInvalidJSON(t *testing.T) {
	s := store.NewMemoryStore()
	seedUser(t, s)

	router := newTestRouter(s)

	res := authenticatedRequest(
		t,
		router,
		http.MethodPost,
		"/api/testUserID/anthologies",
		strings.NewReader("{invalid"),
	)

	requireStatus(t, res, http.StatusBadRequest)
}
