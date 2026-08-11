package api_test

import (
	"reflect"
	"testing"

	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func TestGetManifest(t *testing.T) {
	router := newTestRouter(nil)

	got := getJSON[stremio.Manifest](
		t,
		router,
		"/token/manifest.json",
	)

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
