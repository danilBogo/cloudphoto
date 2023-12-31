package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
)

var CommandUpload = &cobra.Command{
	Use:   "upload",
	Run:   initUpload,
	Short: "Upload photos to bucket",
}

func initUpload(cmd *cobra.Command, _ []string) {
	album, _ := cmd.Flags().GetString(constants.Album)
	path, _ := cmd.Flags().GetString(constants.Path)

	photos := getPhotosFromDirectory(path)

	configManager, err := services.NewConfigManager()
	services.HandleError(err)

	iniConfig, err := configManager.TryGetConfig()
	services.HandleError(err)

	awsConfig := iniConfig.ToAwsConfig()
	awsManager, err := services.NewAwsManager(awsConfig)
	services.HandleError(err)

	uploadPhotos(awsManager, photos, album, iniConfig.Bucket)
}

func uploadPhotos(awsManager *services.AwsManager, photos []string, album, bucket string) {
	wg := sync.WaitGroup{}
	wg.Add(len(photos))
	for _, photo := range photos {
		go func(photo string) {
			defer wg.Done()
			photoKey := services.GetPhotoKey(album, filepath.Base(photo))
			data, err := os.ReadFile(photo)
			if err != nil {
				fmt.Printf("File %v can not be read\n", photo)
				return
			}

			err = awsManager.UploadPhoto(bucket, photoKey, data)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("File %v successfully uploaded with key %v\n", photo, photoKey)
			}
		}(photo)
	}

	wg.Wait()
}

func getPhotosFromDirectory(path string) []string {
	var jpegFiles []string

	files, err := os.ReadDir(path)
	services.HandleError(err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext == ".jpg" || ext == ".jpeg" {
			jpegFiles = append(jpegFiles, filepath.Join(path, file.Name()))
		}
	}

	if len(jpegFiles) == 0 {
		services.HandleError(errors.New(fmt.Sprintf("In directory %v there is no photos", path)))
	}

	return jpegFiles
}
