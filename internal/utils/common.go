package utils

import (
	"cloudphoto/internal/services"
	"errors"
	"fmt"
	"os"
)

func GetIniConfig() *services.IniConfig {
	configManager, err := services.NewConfigManager()
	HandleError(err)

	isValid, iniConfigFromFile, err := configManager.TryGetConfig()
	HandleError(err)

	if !isValid {
		HandleError(errors.New("ini config file is not valid"))
	}

	return iniConfigFromFile
}

func GetAwsManager(iniConfig *services.IniConfig) *services.AwsManager {
	awsConfig := services.AwsConfig{
		AccessKey:   iniConfig.AccessKey,
		SecretKey:   iniConfig.SecretKey,
		Region:      iniConfig.Region,
		EndpointURL: iniConfig.EndpointURL,
	}

	awsManager, err := services.NewAwsManager(awsConfig)
	HandleError(err)

	return awsManager
}

func HandleErrorWithText(err error, text string) {
	if err != nil {
		fmt.Println(text)
		os.Exit(1)
	}
}

func HandleError(err error) {
	HandleErrorWithText(err, err.Error())
}
