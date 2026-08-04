package api_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/config"
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

func newTestRouter(store api.Store) *gin.Engine {

	cfg := config.LoadViper()

	l := config.ConfigureSlog(cfg.Log)
	slog.SetDefault(l)

	l.Info("logger initialised")

	slog.Info(
		"config found and set",
		slog.Group(
			"app",
			"stremio_id", cfg.App.StremioID,
			"version", cfg.App.Version,
			"name", cfg.App.Name,
			"description", cfg.App.Description,
			"logo_url", cfg.App.LogoURL,
			"main_catalog_name", cfg.App.MainCatalogName,
		),
		slog.Group(
			"log",
			"format_json", cfg.Log.FormatJSON,
			"level", cfg.Log.Level,
		),
	)

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

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Setenv("APP_ENV", "test")

	slog.SetLogLoggerLevel(slog.LevelDebug)

	exitCode := m.Run()
	os.Exit(exitCode)
}
