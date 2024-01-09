package commands

import (
	"cloudphoto/internal/constants"
	"cloudphoto/internal/services"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

var CommandMksite = &cobra.Command{
	Use:   "mksite",
	Run:   initMksite,
	Short: "Generation and publication of photo archive web pages",
}

func initMksite(_ *cobra.Command, _ []string) {
	configManager, err := services.NewConfigManager()
	services.HandleError(err)

	iniConfig, err := configManager.TryGetConfig()
	services.HandleError(err)

	awsConfig := iniConfig.ToAwsConfig()
	awsManager, err := services.NewAwsManager(awsConfig)
	services.HandleError(err)

	setReadPublic(awsManager, iniConfig.Bucket)

	configureStaticWebsite(awsManager, iniConfig.Bucket)

	htmlManager, err := services.NewHtmlManager()
	services.HandleError(err)

	count := generateAlbumsHtml(htmlManager, awsManager, iniConfig.Bucket)

	generateIndexHtml(htmlManager, awsManager, count, iniConfig.Bucket)

	generateErrorHtml(htmlManager, awsManager, iniConfig.Bucket)

	fmt.Printf("http://%v.website.yandexcloud.net/\n", iniConfig.Bucket)
}

func setReadPublic(awsManager *services.AwsManager, bucket string) {
	err := awsManager.PutBucketACL(bucket, s3.BucketCannedACLPublicRead)
	services.HandleError(err)
}

func configureStaticWebsite(awsManager *services.AwsManager, bucket string) {
	err := awsManager.ConfigureStaticWebsite(bucket)
	services.HandleError(err)
}

func generateAlbumsHtml(htmlManager *services.HtmlManager, awsManager *services.AwsManager, bucket string) int {
	prefixes, err := awsManager.GetPrefixes(bucket)
	for prefixIndex, prefix := range prefixes {
		awsPhotos, err := awsManager.GetPhotos(bucket, prefix)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		htmlPhotos := make([]services.Photo, len(awsPhotos))
		for awsPhotoIndex, awsPhoto := range awsPhotos {
			url := services.GetPhotoKey(prefix, awsPhoto)
			htmlPhotos[awsPhotoIndex] = services.Photo{Title: awsPhoto, URL: url}
		}

		data, err := htmlManager.GetAlbumHtml(htmlPhotos)
		if err != nil {
			fmt.Println(err.Error())
		}

		err = awsManager.UploadHTML(bucket, services.GetAlbumName(prefixIndex+1), data)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	services.HandleError(err)

	return len(prefixes)
}

func generateIndexHtml(htmlManager *services.HtmlManager, awsManager *services.AwsManager, count int, bucket string) {
	indexHtml, err := htmlManager.GetIndexHtml(count)
	services.HandleError(err)

	err = awsManager.UploadHTML(bucket, constants.IndexHtml, indexHtml)
	services.HandleError(err)
}

func generateErrorHtml(htmlManager *services.HtmlManager, awsManager *services.AwsManager, bucket string) {
	errorHtml, err := htmlManager.GetErrorHtml()
	services.HandleError(err)

	err = awsManager.UploadHTML(bucket, constants.ErrorHtml, errorHtml)
	services.HandleError(err)
}
