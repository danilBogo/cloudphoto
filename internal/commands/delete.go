package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
)

var CommandDelete = &cobra.Command{
	Use:   "delete",
	Run:   initDelete,
	Short: "Delete albums",
}

func initDelete(cmd *cobra.Command, _ []string) {
	album, _ := cmd.Flags().GetString(constants.Album)
	photo, _ := cmd.Flags().GetString(constants.Photo)

	configManager, err := services.NewConfigManager()
	services.HandleError(err)

	iniConfig, err := configManager.TryGetConfig()
	services.HandleError(err)

	awsConfig := iniConfig.ToAwsConfig()
	awsManager, err := services.NewAwsManager(awsConfig)
	services.HandleError(err)

	if len(photo) == 0 {
		err := awsManager.DeletePhotosByPrefix(iniConfig.Bucket, album)
		services.HandleError(err)

		fmt.Printf("Album %v deleted", album)
	} else {
		err := awsManager.DeletePhoto(iniConfig.Bucket, services.GetPhotoKey(album, filepath.Base(photo)))
		services.HandleError(err)

		fmt.Printf("Photo %v deleted", photo)
	}
}
