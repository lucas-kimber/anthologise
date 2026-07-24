package main

import (
	"github.com/lucas-kimber/anthologise/internal/api"
	"github.com/lucas-kimber/anthologise/internal/config"
)

func main() {

	config.LoadViper()
	config.ConfigureSlog()

	r := api.NewRouter()

	if err := r.Run(":7000"); err != nil {
		panic("Fatal Error, couldn't start Gin: " + err.Error())
	}
}
