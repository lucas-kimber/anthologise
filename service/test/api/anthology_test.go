package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestAddAnthology(t *testing.T) {
	const token = "testtoken"

	want := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
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

	body, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("failed to encode anthology: %v", err)
	}

	s := store.NewMemoryStore()
	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodPost,
		"/testtoken/anthologies",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusCreated,
			res.Code,
		)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/testtoken/catalog/series/anthologise",
		nil,
	)
	res = httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
	}

	var catalog stremio.Catalog
	if err := json.NewDecoder(res.Body).Decode(&catalog); err != nil {
		t.Fatalf("failed to decode catalog: %v", err)
	}

	if len(catalog.Metas) != 1 {
		t.Fatalf(
			"received incorrect number of anthologies: want 1, got %d",
			len(catalog.Metas),
		)
	}

	id := catalog.Metas[0].ID

	if !strings.HasPrefix(id, stremio.AnthologyIDPrefix) {
		t.Errorf(
			"received incorrect anthology ID prefix: got %q",
			id,
		)
	}

	want.ID = id

	req = httptest.NewRequest(
		http.MethodGet,
		"/testtoken/meta/series/"+id,
		nil,
	)
	res = httptest.NewRecorder()

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

func TestUpdateAnthology(t *testing.T) {
	const token = "testtoken"

	original := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          "testid",
			Type:        "series",
			Name:        "Test Anthology",
			Poster:      "Test PosterURL",
			Description: "Test Description",
			Genres:      []string{"Test"},
		},
	}

	want := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          "testid",
			Type:        "series",
			Name:        "Updated Anthology",
			Poster:      "Updated PosterURL",
			Description: "Updated Description",
			Genres:      []string{"Updated"},
		},
		Videos: []stremio.Video{
			{
				ID:       "updated_video",
				Title:    "Updated Video",
				Season:   2,
				Episode:  3,
				Released: "Updated Released",
				Overview: "Updated Overview",
			},
		},
	}

	body, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("failed to encode anthology: %v", err)
	}

	s := store.NewMemoryStore()
	s.AddAnthology(token, original)

	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodPut,
		"/testtoken/anthologies",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/testtoken/meta/series/testid",
		nil,
	)
	res = httptest.NewRecorder()

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

func TestUpdateAnthologyNotFound(t *testing.T) {
	const token = "testtoken"

	anthology := stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:   "missing",
			Type: "series",
			Name: "Missing Anthology",
		},
	}

	body, err := json.Marshal(anthology)
	if err != nil {
		t.Fatalf("failed to encode anthology: %v", err)
	}

	s := store.NewMemoryStore()
	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodPut,
		"/testtoken/anthologies",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusNotFound,
			res.Code,
		)
	}
}

func TestAddAnthologyInvalidJSON(t *testing.T) {
	s := store.NewMemoryStore()
	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodPost,
		"/testtoken/anthologies",
		bytes.NewBufferString("{invalid"),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusBadRequest,
			res.Code,
		)
	}
}

func TestUpdateAnthologyInvalidJSON(t *testing.T) {
	s := store.NewMemoryStore()
	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodPut,
		"/testtoken/anthologies",
		bytes.NewBufferString("{invalid"),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusBadRequest,
			res.Code,
		)
	}
}
