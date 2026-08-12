package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/store"
)

func TestAuth(t *testing.T) {
	s := store.NewMemoryStore()
	router := newTestRouter(s)

	anthology := testAnthology("", "Test Anthology")

	// Unauthenticated request should fail.
	res := jsonRequest(
		t,
		router,
		http.MethodPost,
		"/api/testUserID/anthologies",
		anthology,
	)

	requireStatus(t, res, http.StatusUnauthorized)

	// Create a user.
	res = request(
		t,
		router,
		http.MethodPost,
		"/api/users",
		nil,
	)

	requireStatus(t, res, http.StatusCreated)

	user := decodeJSON[struct {
		UserID string `json:"userID"`
	}](t, res)

	if user.UserID == "" {
		t.Fatal("expected userID")
	}

	// Find the edit token cookie.
	var editTokenCookie *http.Cookie

	for _, cookie := range res.Result().Cookies() {
		if cookie.Name == "__Host-anthologise_edit_token" {
			editTokenCookie = cookie
			break
		}
	}

	if editTokenCookie == nil {
		t.Fatal("expected edit token cookie")
	}

	// Use the created user's cookie on an authenticated request.
	body, err := json.Marshal(anthology)
	if err != nil {
		t.Fatalf("failed to encode anthology: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/"+user.UserID+"/anthologies",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(editTokenCookie)

	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	requireStatus(t, res, http.StatusCreated)
}
