package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func TestGetManifest(t *testing.T) {

	router := newTestRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"received incorrect status code: want %d, got %d",
			http.StatusOK,
			res.Code,
		)
	}

	var got stremio.Manifest
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	want := stremio.NewManifest(stremio.ManifestConfig{
		ID:          testID,
		Version:     testVersion,
		Name:        testName,
		Description: testDescription,
		Logo:        testLogo,
		CatalogName: testCatalogName,
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"received incorrect manifest:\nwant: %+v\ngot:  %+v",
			want,
			got,
		)
	}
}
