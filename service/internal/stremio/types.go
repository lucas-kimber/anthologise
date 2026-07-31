// Types for the different meta objects StremIO expects.
// For example, Catalogs and Series.
package stremio

const (
	catalogID         = "anthologise"
	catalogName       = "My Anthologise"
	anthologyIDPrefix = "anthologies_"
	seriesType        = "series"
)

type ManifestConfig struct {
	ID          string
	Version     string
	Name        string
	Description string
	Logo        string
	CatalogName string
}

// Describes the StremIO manifest for anthologise.
// This contains all the plugin info:
// https://stremio.github.io/stremio-addon-guide/step1
type Manifest struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	Resources   []any     `json:"resources"`
	Types       []string  `json:"types"`
	Catalogs    []Catalog `json:"catalogs"`
}

type Resource struct {
	Name       string   `json:"name"`
	Types      []string `json:"types"`
	IDPrefixes []string `json:"idPrefixes,omitempty"`
}

type Catalog struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func NewManifest(cfg ManifestConfig) Manifest {
	return Manifest{
		ID:          cfg.ID,
		Version:     cfg.Version,
		Name:        cfg.Name,
		Description: cfg.Description,
		Logo:        cfg.Logo,
		Resources: []any{
			"catalog",
			Resource{
				Name:       "meta",
				Types:      []string{seriesType},
				IDPrefixes: []string{anthologyIDPrefix},
			},
		},
		Types: []string{
			seriesType,
		},
		Catalogs: []Catalog{
			{
				ID:   catalogID,
				Type: seriesType,
				Name: cfg.CatalogName,
			},
		},
	}
}

type Meta struct {
}

// An Anthology acts as a series wrapper for vidios.
type Anthology struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Poster      string   `json:"poster,omitempty"`
	Description string   `json:"description,omitempty"`
	Genres      []string `json:"genres,omitempty"`

	Videos []Video `json:"videos,omitempty"`
}

// A video is the metadata for a single playable media instance.
// ID is the IMDB for that episode or film.
type Video struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Season   int    `json:"season"`
	Episode  int    `json:"episode"`
	Released string `json:"released,omitempty"`
	Overview string `json:"overview,omitempty"`
}
