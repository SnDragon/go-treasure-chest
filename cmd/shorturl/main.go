package main

import (
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl"
	"log"
)

func main() {
	app := &shorturl.App{}
	if err := app.Initialize(); err != nil {
		log.Fatalf("app init err: %+v\n", err)
	}
	app.Run(":80")
}
