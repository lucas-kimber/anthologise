package store

import "github.com/lucas-kimber/anthologise/internal/stremio"

func GetManifestByToken(token string) (stremio.Manifest, bool) {
	return stremio.Manifest{
		ID:          "dev.anthologise.poc",
		Version:     "0.1.0",
		Name:        "anthologise",
		Description: "Ipsum Lorum",
		Resources: []any{
			"catalog",
			map[string]any{
				"name":       "meta",
				"types":      []string{"series"},
				"idPrefixes": []string{"anthologise_"},
			},
		},
		Types: []string{"series"},
		Catalogs: []stremio.Catalog{
			{
				Type: "series",
				ID:   "anthologise",
				Name: "anthologise",
			},
		},
	}, true
}

func GetCatalogsByToken(token string) ([]stremio.Anthology, bool) {
	return []stremio.Anthology{
		{
			ID:          "anthologise_star_trek_poc",
			Type:        "series",
			Name:        "anthologise - Star Trek POC",
			Poster:      "https://images.metahub.space/poster/medium/tt0244365/img",
			Description: "A deliberately mixed sequence of Enterprise and Discovery episodes.",
			Genres:      []string{"Science Fiction", "Proof of Concept"},
		},
	}, true
}

func GetAnthologiesByToken(token string) (stremio.Anthology, bool) {
	return stremio.Anthology{
		ID:          "anthologise_star_trek_poc",
		Type:        "series",
		Name:        "anthologise - Star Trek POC",
		Poster:      "https://images.metahub.space/poster/medium/tt0244365/img",
		Description: "A deliberately mixed sequence of Enterprise and Discovery episodes.",
		Genres:      []string{"Science Fiction", "Proof of Concept"},
		Videos: []stremio.Video{
			{
				ID:       "tt0244365:1:1",
				Title:    "1. Enterprise S01E01",
				Season:   1,
				Episode:  1,
				Released: "2001-09-26T00:00:00.000Z",
			},
			{
				ID:       "tt0244365:1:2",
				Title:    "2. Enterprise S01E02",
				Season:   1,
				Episode:  2,
				Released: "2001-10-03T00:00:00.000Z",
			},
			{
				ID:       "tt5171438:1:1",
				Title:    "3. Discovery S01E01",
				Season:   1,
				Episode:  3,
				Released: "2017-09-24T00:00:00.000Z",
			},
			{
				ID:       "tt5171438:1:2",
				Title:    "4. Discovery S01E02",
				Season:   1,
				Episode:  4,
				Released: "2017-09-24T00:00:00.000Z",
			},
		},
	}, true
}
