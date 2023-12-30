package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"cloudphoto/internal/utils"
	"fmt"
	"github.com/spf13/cobra"
)

var CommandInit = &cobra.Command{
	Use:   "init",
	Run:   initFunc,
	Short: "Initialize command",
}

func initFunc(_ *cobra.Command, _ []string) {
	bucket := scanValue("Enter bucket")
	accessKey := scanValue("Enter access key id")
	secretKey := scanValue("Enter secret access key")
	iniConfig := &services.IniConfig{
		Bucket:      bucket,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Region:      constants.CurrentRegion,
		EndpointURL: constants.CurrentEndpointURL,
	}

	generateOrUpdateIni(iniConfig)

	createBucketIfNotExist(iniConfig)
}

func scanValue(printingValue string) string {
	var result string
	fmt.Println(printingValue)
	_, err := fmt.Scan(&result)
	utils.HandleError(err)

	return result
}

func generateOrUpdateIni(iniConfig *services.IniConfig) {
	configManager, err := services.NewConfigManager()
	utils.HandleError(err)

	err = configManager.GenerateIni(iniConfig)
	utils.HandleError(err)
}

func createBucketIfNotExist(iniConfig *services.IniConfig) {
	awsManager := utils.GetAwsManager(iniConfig)

	exists, err := awsManager.BucketExists(iniConfig.Bucket)
	utils.HandleErrorWithText(err, fmt.Sprintf("Bucket with name %v already exists", iniConfig.Bucket))

	if !exists {
		err := awsManager.CreateBucket(iniConfig.Bucket)
		utils.HandleErrorWithText(err, fmt.Sprintf("Can not create bucket with name %v", iniConfig.Bucket))

		fmt.Printf("Bucket with name '%v' created\n", iniConfig.Bucket)
	} else {
		fmt.Printf("Bucket with name '%v' exists\n", iniConfig.Bucket)
	}
}
