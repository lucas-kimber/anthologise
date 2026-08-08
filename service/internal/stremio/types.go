// Types for the different protocol types StremIO expects.
// For example, Catalogs and Series.
package stremio

const (
	MainCatalogID     = "anthologise"
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

// Manifest is the StremIO manifest for Anthologise.
// This contains all the plugin info:
// https://stremio.github.io/stremio-addon-guide/step1
type Manifest struct {
	ID          string            `json:"id"`
	Version     string            `json:"version"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Logo        string            `json:"logo"`
	Resources   []Resource        `json:"resources"`
	Types       []string          `json:"types"`
	Catalogs    []ManifestCatalog `json:"catalogs"`
}

type Resource struct {
	Name       string   `json:"name"`
	Types      []string `json:"types"`
	IDPrefixes []string `json:"idPrefixes,omitempty"`
}

// ManifestCatalog describes the catalogs available under a manifest.
type ManifestCatalog struct {
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
		Resources: []Resource{
			{
				Name:  "catalog",
				Types: []string{seriesType},
			},
			{
				Name:       "meta",
				Types:      []string{seriesType},
				IDPrefixes: []string{anthologyIDPrefix},
			},
		},
		Types: []string{
			seriesType,
		},
		Catalogs: []ManifestCatalog{
			{
				ID:   MainCatalogID,
				Type: seriesType,
				Name: cfg.CatalogName,
			},
		},
	}
}

// Catalog represents a collection of metas, in this case Anthologies
type Catalog struct {
	Metas []AnthologyPreview `json:"metas"`
}

// AnthologyPreview represents an anthology within a catalog
type AnthologyPreview struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Poster      string   `json:"poster"`
	Description string   `json:"description,omitempty"`
	Genres      []string `json:"genres,omitempty"`
}

// Anthology acts as a series wrapper for videos, as well as the preview information.
type Anthology struct {
	AnthologyPreview

	Videos []Video `json:"videos"`
}

// Video is the metadata for a single playable media instance.
// ID any ID for that episode or film that Cinemata can look up.
type Video struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Season   int    `json:"season"`
	Episode  int    `json:"episode"`
	Released string `json:"released,omitempty"`
	Overview string `json:"overview,omitempty"`
}
