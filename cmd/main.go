package main

import (
	"cloudphoto/internal/app"
	"cloudphoto/internal/utils"
)

func main() {
	a, err := app.NewApp()
	utils.HandleError(err)

	err = a.AddCommands()
	utils.HandleError(err)
}
