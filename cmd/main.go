package main

import (
	"cloudphoto/internal/app"
	"cloudphoto/internal/services"
)

func main() {
	a, err := app.NewApp()
	services.HandleError(err)

	err = a.AddCommands()
	services.HandleError(err)
}
