package main

import (
	"github.com/lucas-kimber/anylist-poc/internal/api"
)

func main() {

	r := api.NewRouter()

	if err := r.Run(":7000"); err != nil {
		panic("Fatal Error, couldn't start Gin: " + err.Error())
	}
}
