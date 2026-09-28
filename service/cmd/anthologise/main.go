package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucas-kimber/anthologise/service/internal/config"
	"github.com/lucas-kimber/anthologise/service/internal/httpserver"
	"github.com/lucas-kimber/anthologise/service/internal/store"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

const startupTimeout = 10 * time.Second
const shutdownTimeout = 10 * time.Second

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

	startupCtx, cancelStartup := context.WithTimeout(
		context.Background(),
		startupTimeout,
	)
	defer cancelStartup()

	store, err := store.NewPostgresStore(
		startupCtx,
		cfg.DB.DatabaseURL,
	)

	if err != nil {
		slog.Error("database failed to connect", "url", cfg.DB.DatabaseURL, "error", err)
		panic("Fatal error, couldn't connect to database: " + err.Error())
	}

	defer store.Close()

	r := httpserver.NewRouter(manifest, store)

	srv := &http.Server{
		Addr:    ":7000",
		Handler: r,
	}

	srvError := make(chan error, 1)

	go func() {

		err := srv.ListenAndServe()

		slog.Info("HTTP server starting", "address", srv.Addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
			return
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")

	case err := <-srvError:
		slog.Error("HTTP server failer", "error", err)
		return
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)

		if err := srv.Close(); err != nil {
			slog.Error("HTTP server failed to close", "error", err)
		}
	}

	slog.Info("HTTP server stopped")
}
