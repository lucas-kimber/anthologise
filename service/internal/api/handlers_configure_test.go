package api

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
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type handlersConfigTestStore struct {
	Store

	createdUserID string
	createUserErr error

	createdAnthology   stremio.Anthology
	createAnthologyErr error

	setCatalogUserID       string
	setCatalogAnthologyIDs []string
	setCatalogErr          error
}

func (s *handlersConfigTestStore) CreateUser(ctx context.Context, userID string, tokenHash [32]byte) error {

	s.createdUserID = userID

	return s.createUserErr
}

func (s *handlersConfigTestStore) CreateAnthology(ctx context.Context, anthology stremio.Anthology) error {

	s.createdAnthology = anthology

	return s.createAnthologyErr
}

func (s *handlersConfigTestStore) SetCatalog(ctx context.Context, userID string, anthologyIDs []string) error {

	s.setCatalogUserID = userID
	s.setCatalogAnthologyIDs = anthologyIDs

	return s.setCatalogErr
}

func newConfigHandlersTestRouter(store *handlersConfigTestStore) *gin.Engine {

	server := server{
		store: store,
	}

	r := gin.New()

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

			store := &handlersConfigTestStore{
				createUserErr: tt.storeErr,
			}

			r := newConfigHandlersTestRouter(store)

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

			store := &handlersConfigTestStore{
				createAnthologyErr: tt.createAnthologyErr,
			}

			r := newConfigHandlersTestRouter(store)

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

			if !strings.HasPrefix(got.ID, stremio.AnthologyIDPrefix+"_") {
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
			storeErr: ErrAnthologyNotFound,
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

			store := &handlersConfigTestStore{
				setCatalogErr: tt.storeErr,
			}

			r := newConfigHandlersTestRouter(store)

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
