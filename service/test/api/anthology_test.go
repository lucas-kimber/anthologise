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
