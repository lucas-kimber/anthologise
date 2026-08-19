package api_test

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/config"
	"github.com/lucas-kimber/anthologise/service/internal/store"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

const (
	testUserID      = "testUserID"
	testEditToken   = "testEditToken"
	testID          = "testid"
	testVersion     = "testversion"
	testName        = "testname"
	testDescription = "testdescription"
	testLogo        = "testlogo"
	testCatalogName = "testcatalog"
)

func newTestRouter(store api.Store) *gin.Engine {
	cfg := config.LoadViper()

	manifest := stremio.NewManifest(stremio.ManifestConfig{
		ID:          cfg.App.StremioID,
		Version:     cfg.App.Version,
		Name:        cfg.App.Name,
		Description: cfg.App.Description,
		Logo:        cfg.App.LogoURL,
		CatalogName: cfg.App.MainCatalogName,
	})

	return api.NewRouter(manifest, store)
}

func seedUser(t *testing.T, s *store.MemoryStore) {
	t.Helper()

	tokenHash := sha256.Sum256([]byte(testEditToken))

	if err := s.CreateUser(context.Background(), testUserID, tokenHash); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	env := map[string]string{
		"ANTHOLOGISE_STREMIO_ID":           testID,
		"ANTHOLOGISE_VERSION_NUMBER":       testVersion,
		"ANTHOLOGISE_APP_NAME":             testName,
		"ANTHOLOGISE_MANIFEST_DESCRIPTION": testDescription,
		"ANTHOLOGISE_LOGO_URL":             testLogo,
		"ANTHOLOGISE_MAIN_CATALOG_NAME":    testCatalogName,
	}

	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}

	cfg := config.LoadViper()
	slog.SetDefault(config.ConfigureSlog(cfg.Log))

	os.Exit(m.Run())
}
