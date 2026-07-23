package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	addonID   = "dev.anylist.poc"
	listID    = "anylist_star_trek_poc"
	catalogID = "anylist"
)

type Manifest struct {
	ID          string            `json:"id"`
	Version     string            `json:"version"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Resources   []any             `json:"resources"`
	Types       []string          `json:"types"`
	Catalogs    []ManifestCatalog `json:"catalogs"`
	IDPrefixes  []string          `json:"idPrefixes,omitempty"`
}

type ManifestCatalog struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MetaPreview struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Poster      string   `json:"poster,omitempty"`
	Description string   `json:"description,omitempty"`
	Genres      []string `json:"genres,omitempty"`
}

type Video struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Season   int    `json:"season"`
	Episode  int    `json:"episode"`
	Released string `json:"released,omitempty"`
	Overview string `json:"overview,omitempty"`
}

type Meta struct {
	MetaPreview
	Videos []Video `json:"videos"`
}

var manifest = Manifest{
	ID:          addonID,
	Version:     "0.1.0",
	Name:        "AnyList POC",
	Description: "Tests a custom episode order spanning multiple series.",
	Resources: []any{
		"catalog",
		map[string]any{
			"name":       "meta",
			"types":      []string{"series"},
			"idPrefixes": []string{"anylist_"},
		},
	},
	Types: []string{"series"},
	Catalogs: []ManifestCatalog{
		{
			Type: "series",
			ID:   catalogID,
			Name: "AnyList",
		},
	},
}

var starTrekList = Meta{
	MetaPreview: MetaPreview{
		ID:          listID,
		Type:        "series",
		Name:        "AnyList - Star Trek POC",
		Poster:      "https://images.metahub.space/poster/medium/tt0244365/img",
		Description: "A deliberately mixed sequence of Enterprise and Discovery episodes.",
		Genres:      []string{"Science Fiction", "Proof of Concept"},
	},
	Videos: []Video{
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
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "AnyList POC\nInstall: %s/manifest.json\n", requestBaseURL(c))
	})

	router.GET("/manifest.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, manifest)
	})

	router.GET("/catalog/:mediaType/*path", catalogHandler)
	router.GET("/meta/:mediaType/:id", metaHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	log.Printf("AnyList POC listening on http://127.0.0.1:%s", port)
	log.Printf("Install http://127.0.0.1:%s/manifest.json in Stremio", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func catalogHandler(c *gin.Context) {
	if c.Param("mediaType") != "series" {
		c.JSON(http.StatusOK, gin.H{"metas": []MetaPreview{}})
		return
	}

	path := strings.TrimPrefix(c.Param("path"), "/")
	parts := strings.Split(path, "/")
	requestedCatalog := strings.TrimSuffix(parts[0], ".json")

	if requestedCatalog != catalogID {
		c.JSON(http.StatusOK, gin.H{"metas": []MetaPreview{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metas": []MetaPreview{starTrekList.MetaPreview},
	})
}

func metaHandler(c *gin.Context) {
	mediaType := c.Param("mediaType")
	id := strings.TrimSuffix(c.Param("id"), ".json")

	if mediaType != "series" || id != listID {
		c.JSON(http.StatusOK, gin.H{"meta": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": starTrekList})
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func requestBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
