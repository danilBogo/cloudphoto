package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
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
	services.HandleError(err)

	return result
}

func generateOrUpdateIni(iniConfig *services.IniConfig) {
	configManager, err := services.NewConfigManager()
	services.HandleError(err)

	err = configManager.GenerateIni(iniConfig)
	services.HandleError(err)
}

func createBucketIfNotExist(iniConfig *services.IniConfig) {
	awsConfig := iniConfig.ToAwsConfig()
	awsManager, err := services.NewAwsManager(awsConfig)
	services.HandleError(err)

	exists, err := awsManager.BucketExists(iniConfig.Bucket)
	services.HandleErrorWithText(err, fmt.Sprintf("Bucket with name %v already exists", iniConfig.Bucket))

	if !exists {
		err := awsManager.CreateBucket(iniConfig.Bucket)
		services.HandleErrorWithText(err, fmt.Sprintf("Can not create bucket with name %v", iniConfig.Bucket))

		fmt.Printf("Bucket with name '%v' created\n", iniConfig.Bucket)
	} else {
		fmt.Printf("Bucket with name '%v' exists\n", iniConfig.Bucket)
	}
}
