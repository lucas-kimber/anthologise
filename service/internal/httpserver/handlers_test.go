package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/store"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type handlersTestStore struct {
	Store

	createdUserID string
	createUserErr error

	createdAnthology   stremio.Anthology
	createAnthologyErr error

	setCatalogUserID       string
	setCatalogAnthologyIDs []string
	setCatalogErr          error

	userID      string
	anthologyID string
	anthology   stremio.Anthology
	catalog     stremio.Catalog
}

func (s *handlersTestStore) CreateUser(ctx context.Context, userID string, tokenHash [32]byte) error {

	s.createdUserID = userID

	return s.createUserErr
}

func (s *handlersTestStore) CreateAnthology(ctx context.Context, anthology stremio.Anthology) error {

	s.createdAnthology = anthology

	return s.createAnthologyErr
}

func (s *handlersTestStore) SetCatalog(ctx context.Context, userID string, anthologyIDs []string) error {

	s.setCatalogUserID = userID
	s.setCatalogAnthologyIDs = anthologyIDs

	return s.setCatalogErr
}

func (s *handlersTestStore) GetCatalog(ctx context.Context, userID string) (stremio.Catalog, error) {

	if s.userID != userID {
		return stremio.Catalog{}, store.ErrUserDoesNotExist
	}

	return s.catalog, nil
}

func (s *handlersTestStore) GetAnthology(ctx context.Context, anthologyID string) (stremio.Anthology, error) {

	if s.anthologyID != anthologyID {
		return stremio.Anthology{}, store.ErrAnthologyNotFound
	}

	return s.anthology, nil
}

func newHandlersTestRouter(store *handlersTestStore) *gin.Engine {

	server := server{
		manifest: testManifest,
		store:    store,
	}

	r := gin.New()

	r.GET("/:userID/manifest.json", server.getManifest)
	r.GET("/:userID/catalog/:type/:anthologyID", server.getCatalog)
	r.GET("/:userID/meta/:type/:anthologyID", server.getAnthology)

	r.POST("/api/users", server.createUser)
	r.POST("/api/anthologies", server.createAnthology)
	r.PUT("/api/:userID/catalog", server.setCatalog)

	return r
}

func TestCreateUser(t *testing.T) {

	tests := []struct {
		name     string
		storeErr error
		want     int
	}{
		{
			name: "success",
			want: http.StatusCreated,
		},
		{
			name:     "store failure",
			storeErr: errors.New("test store error"),
			want:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			store := &handlersTestStore{
				createUserErr: tt.storeErr,
			}

			r := newHandlersTestRouter(store)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/users",
				nil,
			)

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != tt.want {
				t.Errorf("got %d, want %d", res.Code, tt.want)
			}

			if tt.want != http.StatusCreated {
				return
			}

			var got struct {
				UserID string `json:"userID"`
			}

			if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if got.UserID != store.createdUserID {
				t.Errorf("got userID %s, want %s", got.UserID, store.createdUserID)
			}

			if len(res.Result().Cookies()) == 0 {
				t.Error("expected edit token cookie")
			}
		})
	}
}

func TestCreateAnthology(t *testing.T) {

	body, err := json.Marshal(testAnthology)
	if err != nil {
		t.Fatalf("failed to marshal test anthology: %v", err)
	}

	tests := []struct {
		name               string
		body               string
		createAnthologyErr error
		want               int
	}{
		{
			name: "success",
			body: string(body),
			want: http.StatusCreated,
		},
		{
			name: "invalid anthology",
			body: "{invalid",
			want: http.StatusBadRequest,
		},
		{
			name:               "store failure",
			body:               string(body),
			createAnthologyErr: errors.New("test store error"),
			want:               http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			store := &handlersTestStore{
				createAnthologyErr: tt.createAnthologyErr,
			}

			r := newHandlersTestRouter(store)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/anthologies",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != tt.want {
				t.Errorf("got %d, want %d", res.Code, tt.want)
			}

			if tt.want != http.StatusCreated {
				return
			}

			var got stremio.Anthology

			if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if !reflect.DeepEqual(got, store.createdAnthology) {
				t.Errorf("got %+v, want %+v", got, store.createdAnthology)
			}

			if !strings.HasPrefix(got.ID, stremio.AnthologyIDPrefix) {
				t.Errorf("invalid anthology ID: %s", got.ID)
			}
		})
	}
}

func TestSetCatalog(t *testing.T) {

	tests := []struct {
		name     string
		body     string
		storeErr error
		wantIDs  []string
		want     int
	}{
		{
			name:    "success",
			body:    `{"anthologyIDs":["anthology-1","anthology-2"]}`,
			wantIDs: []string{"anthology-1", "anthology-2"},
			want:    http.StatusNoContent,
		},
		{
			name:    "empty catalog",
			body:    `{"anthologyIDs":[]}`,
			wantIDs: []string{},
			want:    http.StatusNoContent,
		},
		{
			name: "invalid request",
			body: "{invalid",
			want: http.StatusBadRequest,
		},
		{
			name:     "anthology not found",
			body:     `{"anthologyIDs":["anthology-1"]}`,
			storeErr: store.ErrAnthologyNotFound,
			wantIDs:  []string{"anthology-1"},
			want:     http.StatusNotFound,
		},
		{
			name:     "store failure",
			body:     `{"anthologyIDs":["anthology-1"]}`,
			storeErr: errors.New("test store error"),
			wantIDs:  []string{"anthology-1"},
			want:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			store := &handlersTestStore{
				setCatalogErr: tt.storeErr,
			}

			r := newHandlersTestRouter(store)

			req := httptest.NewRequest(
				http.MethodPut,
				"/api/"+testUserID+"/catalog",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != tt.want {
				t.Errorf("got %d, want %d", res.Code, tt.want)
			}

			if tt.want == http.StatusBadRequest {
				return
			}

			if store.setCatalogUserID != testUserID {
				t.Errorf("got userID %s, want %s", store.setCatalogUserID, testUserID)
			}

			if !reflect.DeepEqual(store.setCatalogAnthologyIDs, tt.wantIDs) {
				t.Errorf("got anthologyIDs %v, want %v", store.setCatalogAnthologyIDs, tt.wantIDs)
			}
		})
	}
}

func TestGetManifest(t *testing.T) {

	store := &handlersTestStore{}
	r := newHandlersTestRouter(store)

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

	store := &handlersTestStore{
		userID:  testUserID,
		catalog: testCatalog,
	}

	r := newHandlersTestRouter(store)

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

	store := &handlersTestStore{
		anthologyID: testAnthologyID,
		anthology:   testAnthology,
	}

	r := newHandlersTestRouter(store)

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
