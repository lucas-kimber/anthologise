package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type handlersStremioTestStore struct {
	Store

	userID      string
	anthologyID string
	anthology   stremio.Anthology
	catalog     stremio.Catalog
}

func (s handlersStremioTestStore) GetCatalog(ctx context.Context, userID string) (stremio.Catalog, error) {

	if s.userID != userID {
		return stremio.Catalog{}, ErrUserDoesNotExist
	}

	return s.catalog, nil
}

func (s handlersStremioTestStore) GetAnthology(ctx context.Context, anthologyID string) (stremio.Anthology, error) {

	if s.anthologyID != anthologyID {
		return stremio.Anthology{}, ErrAnthologyNotFound
	}

	return s.anthology, nil
}

func newStremioGettersTestRouter() *gin.Engine {

	store := handlersStremioTestStore{
		userID:      testUserID,
		anthologyID: testAnthologyID,
		anthology:   testAnthology,
		catalog:     testCatalog,
	}

	server := server{
		manifest: testManifest,
		store:    store,
	}

	r := gin.New()

	r.GET("/:userID/manifest.json", server.getManifest)
	r.GET("/:userID/catalog/:type/:anthologyID", server.getCatalog)
	r.GET("/:userID/meta/:type/:anthologyID", server.getAnthology)

	return r
}

func TestGetManifest(t *testing.T) {

	r := newStremioGettersTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/"+testUserID+"/manifest.json",
		nil,
	)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", res.Code, http.StatusOK)
	}

	var got stremio.Manifest

	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(got, testManifest) {
		t.Errorf("got %+v, want %+v", got, testManifest)
	}
}

func TestGetCatalog(t *testing.T) {

	r := newStremioGettersTestRouter()

	tests := []struct {
		name string
		user string
		want int
	}{
		{
			name: "existing user",
			user: testUserID,
			want: http.StatusOK,
		},
		{
			name: "nonexistent user",
			user: "missing",
			want: http.StatusNotFound,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+tt.user+"/catalog/series/"+testAnthologyID,
				nil,
			)

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != tt.want {
				t.Errorf("got %d, want %d", res.Code, tt.want)
			}

			if tt.want != http.StatusOK {
				return
			}

			var got stremio.Catalog

			if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if !reflect.DeepEqual(got, testCatalog) {
				t.Errorf("got %+v, want %+v", got, testCatalog)
			}
		})
	}
}

func TestGetAnthology(t *testing.T) {

	r := newStremioGettersTestRouter()

	tests := []struct {
		name        string
		anthologyID string
		want        int
	}{
		{
			name:        "existing anthology",
			anthologyID: testAnthologyID,
			want:        http.StatusOK,
		},
		{
			name:        "nonexistent anthology",
			anthologyID: "missing",
			want:        http.StatusNotFound,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+testUserID+"/meta/series/"+tt.anthologyID,
				nil,
			)

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != tt.want {
				t.Errorf("got %d, want %d", res.Code, tt.want)
			}

			if tt.want != http.StatusOK {
				return
			}

			var got stremio.Anthology

			if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if !reflect.DeepEqual(got, testAnthology) {
				t.Errorf("got %+v, want %+v", got, testAnthology)
			}
		})
	}
}
