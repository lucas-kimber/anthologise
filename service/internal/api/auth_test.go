package api

import (
	"context"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type authTestStore struct {
	Store

	userID    string
	tokenHash [32]byte
}

func (s authTestStore) GetTokenHash(ctx context.Context, userID string) ([32]byte, error) {

	if userID != s.userID {
		var b [32]byte
		return b, ErrUserDoesNotExist
	}

	return s.tokenHash, nil
}

func TestAuth(t *testing.T) {

	testUserID := "00000000-0000-0000-0000-000000000001"

	store := authTestStore{
		userID:    testUserID,
		tokenHash: sha256.Sum256([]byte(testEditToken)),
	}

	server := server{
		store: store,
	}

	r := gin.New()
	private := r.Group("/api/:userID")
	private.Use(server.authMiddleware())

	private.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		name      string
		user      string
		editToken string
		want      int
	}{
		{
			name:      "missing token",
			user:      testUserID,
			editToken: "",
			want:      http.StatusUnauthorized,
		},
		{
			name:      "invalid token",
			user:      testUserID,
			editToken: "invalid",
			want:      http.StatusUnauthorized,
		},
		{
			name:      "nonexistant user",
			user:      "missing",
			editToken: testEditToken,
			want:      http.StatusUnauthorized,
		},
		{
			name:      "valid token",
			user:      testUserID,
			editToken: testEditToken,
			want:      http.StatusOK,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/"+tt.user+"/test",
				nil,
			)

			if tt.editToken != "" {
				req.AddCookie(&http.Cookie{
					Name:  "__Host-anthologise_edit_token",
					Value: tt.editToken,
				})
			}

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			got := res.Result().StatusCode

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
