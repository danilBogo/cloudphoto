package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"cloudphoto/internal/utils"
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

	iniConfig := utils.GetIniConfig()

	awsManager := utils.GetAwsManager(iniConfig)

	if len(photo) == 0 {
		err := awsManager.DeletePhotosByPrefix(iniConfig.Bucket, album)
		utils.HandleError(err)

		fmt.Printf("Album %v deleted", album)
	} else {
		err := awsManager.DeletePhoto(iniConfig.Bucket, services.GetPhotoKey(album, filepath.Base(photo)))
		utils.HandleError(err)

		fmt.Printf("Photo %v deleted", photo)
	}
}
