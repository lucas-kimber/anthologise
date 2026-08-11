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

const testToken = "testtoken"

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

func jsonRequest(
	t *testing.T,
	handler http.Handler,
	method string,
	path string,
	body any,
) *httptest.ResponseRecorder {
	t.Helper()

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}

	req := httptest.NewRequest(
		method,
		path,
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	return res
}

func decodeJSON[T any](
	t *testing.T,
	res *httptest.ResponseRecorder,
) T {
	t.Helper()

	var got T

	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return got
}

func requireStatus(
	t *testing.T,
	res *httptest.ResponseRecorder,
	want int,
) {
	t.Helper()

	if res.Code != want {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			want,
			res.Code,
		)
	}
}

func addToCatalog(
	t *testing.T,
	s *store.MemoryStore,
	token string,
	anthologies ...stremio.Anthology,
) {
	t.Helper()

	for _, anthology := range anthologies {
		if err := s.CreateAnthology(token, anthology); err != nil {
			t.Fatalf("failed to create anthology: %v", err)
		}

		if err := s.AddAnthologyToCatalog(token, anthology.ID); err != nil {
			t.Fatalf("failed to add anthology to catalog: %v", err)
		}
	}
}

func TestGetAnthology(t *testing.T) {
	want := testAnthology("testid", "Test Anthology")

	s := store.NewMemoryStore()

	if err := s.CreateAnthology(testToken, want); err != nil {
		t.Fatalf("failed to create anthology: %v", err)
	}

	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodGet,
		"/testtoken/meta/series/testid",
		nil,
	)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusOK)

	got := decodeJSON[stremio.Anthology](t, res)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect anthology:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}

func TestGetCatalogWithTwoAnthologies(t *testing.T) {
	first := testAnthology("testid1", "Test Anthology 1")
	second := testAnthology("testid2", "Test Anthology 2")

	s := store.NewMemoryStore()

	addToCatalog(
		t,
		s,
		testToken,
		first,
		second,
	)

	router := newTestRouter(s)

	req := httptest.NewRequest(
		http.MethodGet,
		"/testtoken/catalog/series/anthologise",
		nil,
	)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusOK)

	want := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			first.AnthologyPreview,
			second.AnthologyPreview,
		},
	}

	got := decodeJSON[stremio.Catalog](t, res)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect catalog:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}

func TestAddAnthology(t *testing.T) {
	want := testAnthology("", "Test Anthology")

	s := store.NewMemoryStore()
	router := newTestRouter(s)

	res := jsonRequest(
		t,
		router,
		http.MethodPost,
		"/testtoken/anthologies",
		want,
	)

	requireStatus(t, res, http.StatusCreated)

	created := decodeJSON[stremio.Anthology](t, res)

	if !strings.HasPrefix(created.ID, stremio.AnthologyIDPrefix) {
		t.Errorf(
			"generated anthology ID has incorrect prefix: %q",
			created.ID,
		)
	}

	want.ID = created.ID

	if !reflect.DeepEqual(created, want) {
		t.Errorf(
			"created incorrect anthology:\nwant: %+v\ngot:  %+v",
			want,
			created,
		)
	}

	// Check that the anthology was actually stored.

	req := httptest.NewRequest(
		http.MethodGet,
		"/testtoken/meta/series/"+created.ID,
		nil,
	)
	res = httptest.NewRecorder()

	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusOK)

	got := decodeJSON[stremio.Anthology](t, res)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"stored incorrect anthology:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}

	// Creating an anthology should also add it to the creator's catalog.

	req = httptest.NewRequest(
		http.MethodGet,
		"/testtoken/catalog/series/anthologise",
		nil,
	)
	res = httptest.NewRecorder()

	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusOK)

	wantCatalog := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			want.AnthologyPreview,
		},
	}

	gotCatalog := decodeJSON[stremio.Catalog](t, res)

	if !reflect.DeepEqual(gotCatalog, wantCatalog) {
		t.Errorf(
			"received incorrect catalog:\nwant: %+v\ngot:  %+v",
			wantCatalog,
			gotCatalog,
		)
	}
}

func TestUpdateAnthology(t *testing.T) {
	original := testAnthology(
		"testid",
		"Original Anthology",
	)

	want := testAnthology(
		"testid",
		"Updated Anthology",
	)
	want.Description = "Updated Description"

	s := store.NewMemoryStore()

	addToCatalog(
		t,
		s,
		testToken,
		original,
	)

	router := newTestRouter(s)

	res := jsonRequest(
		t,
		router,
		http.MethodPut,
		"/testtoken/anthologies",
		want,
	)

	requireStatus(t, res, http.StatusOK)

	// Check the anthology itself was replaced.

	req := httptest.NewRequest(
		http.MethodGet,
		"/testtoken/meta/series/testid",
		nil,
	)
	res = httptest.NewRecorder()

	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusOK)

	got := decodeJSON[stremio.Anthology](t, res)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect anthology:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}

	// The catalog should resolve the updated preview too.

	req = httptest.NewRequest(
		http.MethodGet,
		"/testtoken/catalog/series/anthologise",
		nil,
	)
	res = httptest.NewRecorder()

	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusOK)

	wantCatalog := stremio.Catalog{
		Metas: []stremio.AnthologyPreview{
			want.AnthologyPreview,
		},
	}

	gotCatalog := decodeJSON[stremio.Catalog](t, res)

	if !reflect.DeepEqual(gotCatalog, wantCatalog) {
		t.Errorf(
			"received incorrect catalog:\nwant: %+v\ngot:  %+v",
			wantCatalog,
			gotCatalog,
		)
	}
}

func TestUpdateAnthologyNotFound(t *testing.T) {
	anthology := testAnthology(
		"missing",
		"Missing Anthology",
	)

	s := store.NewMemoryStore()
	router := newTestRouter(s)

	res := jsonRequest(
		t,
		router,
		http.MethodPut,
		"/testtoken/anthologies",
		anthology,
	)

	requireStatus(t, res, http.StatusNotFound)
}

func TestAnthologyInvalidJSON(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "add",
			method: http.MethodPost,
			path:   "/testtoken/anthologies",
		},
		{
			name:   "update",
			method: http.MethodPut,
			path:   "/testtoken/anthologies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			router := newTestRouter(s)

			req := httptest.NewRequest(
				tt.method,
				tt.path,
				bytes.NewBufferString("{invalid"),
			)
			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			requireStatus(
				t,
				res,
				http.StatusBadRequest,
			)
		})
	}
}
