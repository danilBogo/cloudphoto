package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"github.com/spf13/cobra"
)

var CommandDownload = &cobra.Command{
	Use:   "download",
	Run:   initDownload,
	Short: "Download photos from bucket",
}

func initDownload(cmd *cobra.Command, _ []string) {
	album, _ := cmd.Flags().GetString(constants.Album)
	path, _ := cmd.Flags().GetString(constants.Path)

	configManager, err := services.NewConfigManager()
	services.HandleError(err)

	iniConfig, err := configManager.TryGetConfig()
	services.HandleError(err)

	awsConfig := iniConfig.ToAwsConfig()
	awsManager, err := services.NewAwsManager(awsConfig)
	services.HandleError(err)

	err = awsManager.DownloadPhotos(iniConfig.Bucket, album, path)
	services.HandleError(err)
}
