package app

import (
	commands2 "cloudphoto/internal/commands"
	"cloudphoto/internal/constants"
	"github.com/spf13/cobra"
	"os"
)

type App struct {
	rootCmd    *cobra.Command
	currentDir string
}

func NewApp() (*App, error) {
	currentDir, err := getCurrentDir()
	if err != nil {
		return nil, err
	}

	a := &App{
		rootCmd:    &cobra.Command{Use: "cloudphoto"},
		currentDir: currentDir,
	}

	return a, nil
}

func (a *App) AddCommands() error {
	addInit(a.rootCmd)

	err := addUpload(a.rootCmd, a.currentDir)
	if err != nil {
		return err
	}

	err = addDownload(a.rootCmd, a.currentDir)
	if err != nil {
		return err
	}

	addList(a.rootCmd)

	err = addDelete(a.rootCmd)
	if err != nil {
		return err
	}

	addMksite(a.rootCmd)

	return nil
}

func addInit(rootCmd *cobra.Command) {
	rootCmd.AddCommand(commands2.CommandInit)
}

func addUpload(rootCmd *cobra.Command, currentDir string) error {
	rootCmd.AddCommand(commands2.CommandUpload)
	commands2.CommandUpload.Flags().String(constants.Album, "", "Album name")
	err := commands2.CommandUpload.MarkFlagRequired(constants.Album)
	commands2.CommandUpload.Flags().String(constants.Path, currentDir, "Path to directory with photos")

	return err
}

func addDownload(rootCmd *cobra.Command, currentDir string) error {
	rootCmd.AddCommand(commands2.CommandDownload)
	commands2.CommandDownload.Flags().String(constants.Album, "", "Album name")
	err := commands2.CommandDownload.MarkFlagRequired(constants.Album)
	commands2.CommandDownload.Flags().String(constants.Path, currentDir, "Path to directory with photos")

	return err
}

func addList(rootCmd *cobra.Command) {
	rootCmd.AddCommand(commands2.CommandList)
	commands2.CommandList.Flags().String(constants.Album, "", "Album name")
}

func addDelete(rootCmd *cobra.Command) error {
	rootCmd.AddCommand(commands2.CommandDelete)
	commands2.CommandDelete.Flags().String(constants.Album, "", "Album name")
	err := commands2.CommandDelete.MarkFlagRequired(constants.Album)
	commands2.CommandDelete.Flags().String(constants.Photo, "", "Photo name to delete")

	return err
}

func addMksite(rootCmd *cobra.Command) {
	rootCmd.AddCommand(commands2.CommandMksite)
}

func getCurrentDir() (string, error) {
	return os.Getwd()
}
