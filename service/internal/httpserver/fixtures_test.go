package httpserver

import "github.com/lucas-kimber/anthologise/service/internal/stremio"

const (
	testUserID      = "00000000-0000-0000-0000-000000000001"
	testAnthologyID = "test-anthology-id"
	testEditToken   = "test-edit-token"
)

var (
	testAnthology = stremio.Anthology{
		AnthologyPreview: stremio.AnthologyPreview{
			ID:          testAnthologyID,
			Type:        "series",
			Name:        "Test Name",
			Poster:      "Test PosterURL",
			Description: "Test Description",
			Genres:      []string{"Test"},
		},
		Videos: []stremio.Video{
			{
				ID:       "test_video",
				Title:    "Test Video",
				Season:   1,
				Episode:  1,
				Released: "Test Released",
				Overview: "Test Overview",
			},
		},
	}
	testCatalog = stremio.Catalog{
		Metas: []stremio.AnthologyPreview{testAnthology.AnthologyPreview},
	}
	testManifest = stremio.Manifest{
		ID:          "test-id",
		Version:     "1.0.0",
		Name:        "Test Anthologise",
		Description: "Test Description",
		Logo:        "https://example.com/logo.png",

		Resources: []stremio.Resource{
			{
				Name:  "catalog",
				Types: []string{"series"},
			},
			{
				Name:       "meta",
				Types:      []string{"series"},
				IDPrefixes: []string{stremio.AnthologyIDPrefix},
			},
		},

		Types: []string{
			"series",
		},

		Catalogs: []stremio.ManifestCatalog{
			{
				ID:   stremio.MainCatalogID,
				Type: "series",
				Name: "Test Catalog",
			},
		},
	}
)
