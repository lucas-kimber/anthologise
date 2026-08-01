package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/app"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

const (
	testID          = "testid"
	testVersion     = "testversion"
	testName        = "testname"
	testDescription = "testdescription"
	testLogo        = "testlogo"
	testCatalogName = "testcatalog"
)

func TestGetManifest(t *testing.T) {

	t.Setenv("ANTHOLOGISE_STREMIO_ID", "testid")
	t.Setenv("ANTHOLOGISE_VERSION_NUMBER", "testversion")
	t.Setenv("ANTHOLOGISE_APP_NAME", "testname")
	t.Setenv("ANTHOLOGISE_MANIFEST_DESCRIPTION", "testdescription")
	t.Setenv("ANTHOLOGISE_LOGO_URL", "testlogo")
	t.Setenv("ANTHOLOGISE_MAIN_CATALOG_NAME", "testcatalog")

	router := app.NewApp()

	req, _ := http.NewRequest(http.MethodGet, "/manifest.json", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("go bad response code: want %d, got %d", http.StatusOK, res.Code)
	}

	var got any
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	wantJSON, _ := json.Marshal(stremio.NewManifest(stremio.ManifestConfig{
		ID:          testID,
		Version:     testVersion,
		Name:        testName,
		Description: testDescription,
		Logo:        testLogo,
		CatalogName: testCatalogName,
	}))

	var want any
	_ = json.Unmarshal(wantJSON, &want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("received incorrect manifest from response:\nwant %+v\ngot: %+v", want, got)
	}
}
