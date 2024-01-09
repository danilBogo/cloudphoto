package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"fmt"
	"github.com/spf13/cobra"
)

var CommandList = &cobra.Command{
	Use:   "list",
	Run:   initList,
	Short: "View the list of albums and photos in the album",
}

func initList(cmd *cobra.Command, _ []string) {
	album, _ := cmd.Flags().GetString(constants.Album)

	configManager, err := services.NewConfigManager()
	services.HandleError(err)

	iniConfig, err := configManager.TryGetConfig()
	services.HandleError(err)

	awsConfig := iniConfig.ToAwsConfig()
	awsManager, err := services.NewAwsManager(awsConfig)
	services.HandleError(err)

	var result []string
	if len(album) == 0 {
		result, err = awsManager.GetPrefixes(iniConfig.Bucket)
		services.HandleError(err)
	} else {
		result, err = awsManager.GetPhotos(iniConfig.Bucket, album)
		services.HandleError(err)
	}

	for _, val := range result {
		fmt.Println(val)
	}
}
