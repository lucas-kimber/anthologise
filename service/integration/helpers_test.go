package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

const editTokenCookieName = "__Host-anthologise_edit_token"

type testUser struct {
	ID     string
	Cookie *http.Cookie
}

type createUserResponse struct {
	UserID string `json:"userID"`
}

func loadAnthologyFixture(t *testing.T, name string) stremio.Anthology {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}

	var anthology stremio.Anthology
	if err := json.Unmarshal(data, &anthology); err != nil {
		t.Fatal(err)
	}

	return anthology
}

func doJSONRequest(t *testing.T, method, url string, body any, cookie *http.Cookie) *http.Response {
	t.Helper()

	var reader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatal(err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func requireStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	if resp.StatusCode == expected {
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		t.Fatalf("%s %s: expected status %d, got %d", resp.Request.Method, resp.Request.URL.String(), expected, resp.StatusCode)
	}

	t.Fatalf("%s %s: expected status %d, got %d: %s", resp.Request.Method, resp.Request.URL.String(), expected, resp.StatusCode, body)
}

func decodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()

	var value T

	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}

	return value
}

func findCookie(t *testing.T, resp *http.Response, name string) *http.Cookie {
	t.Helper()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}

	t.Fatalf("response contained no %q cookie", name)
	return nil
}

func createUser(t *testing.T) testUser {
	t.Helper()

	resp := doJSONRequest(t, http.MethodPost, baseURL+"/api/users", nil, nil)
	requireStatus(t, resp, http.StatusCreated)

	cookie := findCookie(t, resp, editTokenCookieName)
	result := decodeJSON[createUserResponse](t, resp)

	if result.UserID == "" {
		t.Fatal("create user response contained empty userID")
	}

	return testUser{ID: result.UserID, Cookie: cookie}
}

func createAnthology(t *testing.T, anthology stremio.Anthology) stremio.Anthology {
	t.Helper()

	resp := doJSONRequest(t, http.MethodPost, baseURL+"/api/anthologies", anthology, nil)
	requireStatus(t, resp, http.StatusCreated)

	return decodeJSON[stremio.Anthology](t, resp)
}

func getAnthology(t *testing.T, userID, anthologyType, anthologyID string) stremio.Anthology {
	t.Helper()

	url := fmt.Sprintf("%s/%s/meta/%s/%s", baseURL, userID, anthologyType, anthologyID)
	resp := doJSONRequest(t, http.MethodGet, url, nil, nil)
	requireStatus(t, resp, http.StatusOK)

	return decodeJSON[stremio.Anthology](t, resp)
}

func getCatalog(t *testing.T, userID string) stremio.Catalog {
	t.Helper()

	url := fmt.Sprintf("%s/%s/catalog/series/anthologise", baseURL, userID)
	resp := doJSONRequest(t, http.MethodGet, url, nil, nil)
	requireStatus(t, resp, http.StatusOK)

	return decodeJSON[stremio.Catalog](t, resp)
}

func requestSetCatalog(t *testing.T, user testUser, anthologyIDs []string) *http.Response {
	t.Helper()

	body := struct {
		AnthologyIDs []string `json:"anthologyIDs"`
	}{
		AnthologyIDs: anthologyIDs,
	}

	url := fmt.Sprintf("%s/api/%s/catalog", baseURL, user.ID)
	return doJSONRequest(t, http.MethodPut, url, body, user.Cookie)
}

func setCatalog(t *testing.T, user testUser, anthologyIDs []string) {
	t.Helper()

	resp := requestSetCatalog(t, user, anthologyIDs)
	defer resp.Body.Close()

	requireStatus(t, resp, http.StatusNoContent)
}
