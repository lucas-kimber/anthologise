package main

import "github.com/lucas-kimber/anthologise/service/internal/app"

func main() {

	r := app.New()

	if err := r.Run(":7000"); err != nil {
		panic("Fatal error, couldn't start Gin: " + err.Error())
	}
}
