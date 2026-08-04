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

func TestGetCatalog(t *testing.T) {

	const token = "testtoken"

	want := stremio.Catalog{
		Metas: []stremio.Anthology{
			{
				AnthologyPreview: stremio.AnthologyPreview{
					ID:          "anthologies_test",
					Type:        "series",
					Name:        "Test Anthology",
					Poster:      "Test PosterURL",
					Description: "Test Description",
					Genres:      []string{"Test"},
				},
			},
		},
	}

	store := store.NewMemoryStore()
	store.AddCatalog(token, stremio.MainCatalogID, want)

	router := newTestRouter(store)

	req := httptest.NewRequest(http.MethodGet, "/testtoken/catalog/series/anthologise.json", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
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
