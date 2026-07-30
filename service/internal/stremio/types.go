// Types for the different meta objects StremIO expects.
// For example, Catalogs and Series.
package stremio

const (
	addonID   = "dev.anthologise.poc"
	listID    = "anthologise_star_trek_poc"
	catalogID = "anthologise"
)

// Describes the StremIO manifest for anthologise.
// This contains all the plugin info:
// https://stremio.github.io/stremio-addon-guide/step1
type Manifest struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resources   []any     `json:"resources"`
	Types       []string  `json:"types"`
	Catalogs    []Catalog `json:"catalogs"`
	IDPrefixes  []string  `json:"idPrefixes,omitempty"`
}

type Meta struct {
	
}

// Describes a single Stremio catalog.
// A catalog is just a shelf associated with the plugin:
// https://stremio.github.io/stremio-addon-guide/step3
type Catalog struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
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
