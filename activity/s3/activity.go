package s3

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

const (
	ovResult = "result"
)

type Activity struct {
	settings   *Settings
	awsSession *session.Session
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var logger log.Logger

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}

	if s.AWSRegion != "" {
		region, err := getRegion(s.AWSRegion)
		if err != nil {
			return nil, err
		}
		act.awsSession = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(region),
		}))
	} else {
		act.awsSession = session.Must(session.NewSession(&aws.Config{}))
	}

	return act, nil
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	in := &Input{}
	err = ctx.GetInputObject(in)
	if err != nil {
		return false, err
	}

	logger = ctx.Logger()
	var action string
	if a.settings.Action != "" {
		action = a.settings.Action
	} else {
		action = in.Action
	}

	var s3err error
	switch action {
	case "download":
		s3err = downloadFileFromS3(a.awsSession, in.LocalLocation, in.S3Location, in.S3BucketName)
	case "upload":
		s3err = uploadFileToS3(a.awsSession, in.LocalLocation, in.S3Location, in.S3BucketName)
	case "delete":
		s3err = deleteFileFromS3(a.awsSession, in.S3Location, in.S3BucketName)
	case "copy":
		s3err = copyFileOnS3(a.awsSession, in.S3Location, in.S3BucketName, in.S3NewLocation)
	case "":
		s3err = errors.New("Action not specified.")
	}
	if s3err != nil {
		// Set the output value in the context
		ctx.SetOutput(ovResult, s3err.Error())
		return true, s3err
	}

	// Set the output value in the context
	ctx.SetOutput(ovResult, "OK")

	return true, nil
}

// Function to download a file from an S3 bucket
func downloadFileFromS3(awsSession *session.Session, directory string, s3Location string, s3BucketName string) error {
	// Create an instance of the S3 Manager
	s3Downloader := s3manager.NewDownloader(awsSession)

	// Create a new temporary file
	f, err := os.Create(filepath.Join(directory, s3Location))
	if err != nil {
		return err
	}

	// Prepare the download
	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3Location),
	}

	// Download the file to disk
	_, err = s3Downloader.Download(f, objectInput)
	if err != nil {
		return err
	}

	return nil
}

// Function to delete a file from an S3 bucket
func deleteFileFromS3(awsSession *session.Session, s3Location string, s3BucketName string) error {
	// Create an instance of the S3 Manager
	s3Session := s3.New(awsSession)

	objectDelete := &s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3Location),
	}

	// Delete the file from S3
	_, err := s3Session.DeleteObject(objectDelete)
	if err != nil {
		return err
	}

	return nil
}

// Function to upload a file from an S3 bucket
func uploadFileToS3(awsSession *session.Session, localFile string, s3Location string, s3BucketName string) error {
	// Create an instance of the S3 Manager
	s3Uploader := s3manager.NewUploader(awsSession)

	// Create a file pointer to the source
	reader, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Prepare the upload
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3Location),
		Body:   reader,
	}

	// Upload the file
	_, err = s3Uploader.Upload(uploadInput)
	if err != nil {
		if reqerr, ok := err.(awserr.RequestFailure); ok {
			logger.Debug("Request failed", reqerr.Code(), reqerr.Message(), reqerr.RequestID())
		} else {
			logger.Debug("Error:", err.Error())
		}
		return err
	}

	return nil
}

// Function to copy a file in an S3 bucket
func copyFileOnS3(awsSession *session.Session, s3Location string, s3BucketName string, s3NewLocation string) error {
	// Create an instance of the S3 Session
	s3Session := s3.New(awsSession)

	// Prepare the copy object
	objectInput := &s3.CopyObjectInput{
		Bucket:     aws.String(s3BucketName),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", s3BucketName, s3Location)),
		Key:        aws.String(s3NewLocation),
	}

	// Copy the object
	_, err := s3Session.CopyObject(objectInput)
	if err != nil {
		return err
	}

	return nil
}

func getRegion(regionSetting string) (string, error) {

	var awsRegions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "cn-north-1", "cn-northwest-1", "eu-central-1", "eu-north-1", "eu-west-1", "eu-west-2", "eu-west-2", "sa-east-1"}

	region := strings.ToLower(regionSetting)
	valid := false
	for _, aRegion := range awsRegions {
		if region == aRegion {
			valid = true
			break
		}
	}

	if !valid {
		return "", fmt.Errorf("unsupported region: %s", regionSetting)
	}

	return region, nil
}
