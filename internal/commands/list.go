package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/utils"
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

	iniConfig := utils.GetIniConfig()

	awsManager := utils.GetAwsManager(iniConfig)

	var result []string
	var err error
	if len(album) == 0 {
		result, err = awsManager.GetPrefixes(iniConfig.Bucket)
		utils.HandleError(err)
	} else {
		result, err = awsManager.GetPhotos(iniConfig.Bucket, album)
		utils.HandleError(err)
	}

	for _, val := range result {
		fmt.Println(val)
	}
}
