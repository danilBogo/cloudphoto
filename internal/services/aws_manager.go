package services

import (
	"bytes"
	"cloudphoto/internal/constants"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type AwsManager struct {
	svc *s3.S3
}

type AwsConfig struct {
	AccessKey   string
	SecretKey   string
	Region      string
	EndpointURL string
}

func NewAwsManager(config AwsConfig) (*AwsManager, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKey,
			config.SecretKey,
			"",
		),
		Endpoint:         aws.String(config.EndpointURL),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	return &AwsManager{svc: s3.New(sess)}, nil
}

func (am AwsManager) BucketExists(bucketName string) (bool, error) {
	_, err := am.svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		var awsError awserr.Error
		if errors.As(err, &awsError) && awsError.Code() == "NotFound" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (am AwsManager) CreateBucket(bucketName string) error {
	_, err := am.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	return err
}

func (am AwsManager) UploadPhoto(bucketName, photoKey string, data []byte) error {
	return am.upload(bucketName, photoKey, aws.String("image/jpeg"), data)
}

func (am AwsManager) UploadHTML(bucketName, htmlKey string, data []byte) error {
	return am.upload(bucketName, htmlKey, aws.String("text/html"), data)
}

func (am AwsManager) DownloadPhotos(bucketName, prefix, path string) error {
	objects, err := am.listObjectsWithPrefix(bucketName, prefix)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(objects))
	for _, object := range objects {
		go func(object *s3.Object) {
			defer wg.Done()

			nameWithoutPrefix := strings.Replace(filepath.Base(*object.Key), prefix+constants.PrefixDivider, "", 1)
			filePath := filepath.Join(path, nameWithoutPrefix)
			err := am.downloadFile(bucketName, *object.Key, filePath)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("File %v successfully uploaded to %v\n", filePath, path)
			}
		}(object)
	}

	wg.Wait()

	return nil
}

func (am AwsManager) GetPrefixes(bucketName string) ([]string, error) {
	objects, err := am.listObjects(bucketName)
	if err != nil {
		return nil, err
	}

	prefixes := getUniquePrefixes(objects)
	if len(prefixes) == 0 {
		return nil, errors.New("no photos in bucket")
	}

	var result []string
	for key, _ := range prefixes {
		result = append(result, key)
	}

	return result, nil
}

func (am AwsManager) GetPhotos(bucketName, prefix string) ([]string, error) {
	objects, err := am.listObjectsWithPrefix(bucketName, prefix)
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, errors.New("no photos with prefix " + prefix)
	}

	var result []string
	for _, object := range objects {
		nameWithoutPrefix := strings.Replace(filepath.Base(*object.Key), prefix+constants.PrefixDivider, "", 1)
		result = append(result, nameWithoutPrefix)
	}

	return result, nil
}

func (am AwsManager) DeletePhoto(bucketName, photoKey string) error {
	_, err := am.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(photoKey),
	})

	if err != nil {
		var s3err awserr.Error
		if errors.As(err, &s3err) && s3err.Code() == "NotFound" {
			return errors.New("no photo to delete")
		}

		return err
	}

	_, err = am.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(photoKey),
	})

	return err
}

func (am AwsManager) DeletePhotosByPrefix(bucketName, prefix string) error {
	objects, err := am.listObjectsWithPrefix(bucketName, prefix)
	if err != nil {
		return err
	}

	objectsToDelete := make([]*s3.ObjectIdentifier, len(objects))
	for i, obj := range objects {
		objectsToDelete[i] = &s3.ObjectIdentifier{
			Key: obj.Key,
		}
	}

	if len(objectsToDelete) == 0 {
		return errors.New("no photos to delete")
	}

	_, err = am.svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: objectsToDelete,
			Quiet:   aws.Bool(false),
		},
	})
	return err
}

func (am AwsManager) PutBucketACL(bucketName, acl string) error {
	_, err := am.svc.PutBucketAcl(&s3.PutBucketAclInput{
		Bucket: aws.String(bucketName),
		ACL:    &acl,
	})

	return err
}

func (am AwsManager) ConfigureStaticWebsite(bucketName string) error {
	_, err := am.svc.PutBucketWebsite(&s3.PutBucketWebsiteInput{
		Bucket: aws.String(bucketName),
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			IndexDocument: &s3.IndexDocument{
				Suffix: aws.String("index.html"),
			},
			ErrorDocument: &s3.ErrorDocument{
				Key: aws.String("error.html"),
			},
		},
	})

	return err
}

func GetPhotoKey(prefix, name string) string {
	return prefix + constants.PrefixDivider + name
}

func (am AwsManager) upload(bucketName, key string, contentType *string, data []byte) error {
	_, err := am.svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        aws.ReadSeekCloser(bytes.NewReader(data)),
		ContentType: contentType,
	})

	return err
}

func (am AwsManager) listObjectsWithPrefix(bucketName, prefix string) ([]*s3.Object, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix + constants.PrefixDivider),
	}

	output, err := am.svc.ListObjects(input)
	if err != nil {
		return nil, err
	}

	if len(output.Contents) == 0 {
		return nil, errors.New("no photos with prefix " + prefix)
	}

	return output.Contents, nil
}

func (am AwsManager) listObjects(bucketName string) ([]*s3.Object, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}

	output, err := am.svc.ListObjects(input)
	if err != nil {
		return nil, err
	}

	if len(output.Contents) == 0 {
		return nil, errors.New("no photos")
	}

	return output.Contents, nil
}

func (am AwsManager) downloadFile(bucketName, objectKey, filePath string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	output, err := am.svc.GetObject(input)
	if err != nil {
		return err
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getUniquePrefixes(objects []*s3.Object) map[string]struct{} {
	result := make(map[string]struct{})
	for _, object := range objects {
		fileName := filepath.Base(*object.Key)
		if !strings.Contains(fileName, "_") {
			continue
		}

		splitStrings := strings.Split(fileName, constants.PrefixDivider)
		if len(splitStrings) == 0 {
			continue
		}

		result[splitStrings[0]] = struct{}{}
	}

	return result
}
