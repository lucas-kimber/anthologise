package api

import "github.com/lucas-kimber/anthologise/service/internal/stremio"

type server struct {
	manifest stremio.Manifest
	store    Store
}

func newServer(manifest stremio.Manifest, store Store) *server {
	return &server{
		manifest: manifest,
		store:    store,
	}
}
