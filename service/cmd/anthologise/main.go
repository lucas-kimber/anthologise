package main

import (
	"log/slog"

	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/config"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func main() {
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

	var store api.Store

	r := api.NewRouter(manifest, store)

	if err := r.Run(":7000"); err != nil {
		panic("Fatal error, couldn't start Gin: " + err.Error())
	}
}
