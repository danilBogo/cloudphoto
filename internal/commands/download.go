package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/utils"
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

	iniConfig := utils.GetIniConfig()

	awsManager := utils.GetAwsManager(iniConfig)

	err := awsManager.DownloadPhotos(iniConfig.Bucket, album, path)
	utils.HandleError(err)
}
