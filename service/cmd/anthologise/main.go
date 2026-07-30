package main

import (
	"github.com/lucas-kimber/anthologise/service/internal/api"
	"github.com/lucas-kimber/anthologise/service/internal/config"
)

func main() {

	cfg := config.LoadViper()
	config.ConfigureSlog(cfg)

	r := api.NewRouter()

	if err := r.Run(":7000"); err != nil {
		panic("Fatal error, couldn't start Gin: " + err.Error())
	}
}
