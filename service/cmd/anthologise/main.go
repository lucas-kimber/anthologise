package main

import (
	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/config"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

func main() {

	cfg := config.LoadViper()
	config.ConfigureSlog(cfg)

	manifest := stremio.NewManifest(stremio.ManifestConfig{
		ID:          cfg.App.StremioID,
		Version:     cfg.App.Version,
		Name:        cfg.App.Name,
		Description: cfg.App.Description,
		Logo:        cfg.App.LogoURL,
	})

	var store api.Store

	r := api.NewRouter(manifest, store)

	if err := r.Run(":7000"); err != nil {
		panic("Fatal error, couldn't start Gin: " + err.Error())
	}
}
