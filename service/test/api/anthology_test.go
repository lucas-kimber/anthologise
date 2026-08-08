package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/store"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func TestGetAnthology(t *testing.T) {

	const token = "testtoken"

	want := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          "testid",
			Type:        "series",
			Name:        "Test Anthology",
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

	s := store.NewMemoryStore()
	s.AddAnthology(token, want)

	router := newTestRouter(s)

	req := httptest.NewRequest(http.MethodGet, "/testtoken/meta/series/testid", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
	}

	var got stremio.Anthology
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect anthology:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}

func TestAddAnthologyUpdatesCatalog(t *testing.T) {
	const token = "testtoken"

	anthology := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          "testid",
			Type:        "series",
			Name:        "Test Anthology",
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

	s := store.NewMemoryStore()
	s.AddAnthology(token, anthology)

	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodGet,
		"/testtoken/catalog/series/anthologise",
		nil,
	)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
	}

	want := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			anthology.AnthologyPreview,
		},
	}

	var got stremio.Catalog
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect catalog:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}

func TestAddTwoAnthologiesUpdatesCatalog(t *testing.T) {
	const token = "testtoken"

	firstAnthology := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          "testid1",
			Type:        "series",
			Name:        "Test Anthology 1",
			Poster:      "Test PosterURL 1",
			Description: "Test Description 1",
			Genres:      []string{"Test"},
		},
		Videos: []stremio.Video{
			{
				ID:       "test_video_1",
				Title:    "Test Video 1",
				Season:   1,
				Episode:  1,
				Released: "Test Released 1",
				Overview: "Test Overview 1",
			},
		},
	}

	secondAnthology := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          "testid2",
			Type:        "series",
			Name:        "Test Anthology 2",
			Poster:      "Test PosterURL 2",
			Description: "Test Description 2",
			Genres:      []string{"Test"},
		},
		Videos: []stremio.Video{
			{
				ID:       "test_video_2",
				Title:    "Test Video 2",
				Season:   1,
				Episode:  1,
				Released: "Test Released 2",
				Overview: "Test Overview 2",
			},
		},
	}

	s := store.NewMemoryStore()
	s.AddAnthology(token, firstAnthology)
	s.AddAnthology(token, secondAnthology)

	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodGet,
		"/testtoken/catalog/series/anthologise",
		nil,
	)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
	}

	want := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			firstAnthology.AnthologyPreview,
			secondAnthology.AnthologyPreview,
		},
	}

	var got stremio.Catalog
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect catalog:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}
