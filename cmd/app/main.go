package main

import (
	"github.com/Blxssy/AvitoTest/config"
	"github.com/Blxssy/AvitoTest/internal/app"
)

func main() {
	app.Run(config.Get())
}
