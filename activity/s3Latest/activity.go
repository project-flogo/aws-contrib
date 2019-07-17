package s3newestmodel

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"io/ioutil"

	"fmt"
	"os"
	"time"

	// "strings"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	// ctx.Logger().Debugf("Setting: %s, %s", s.ReplaceFile, s.UnZip)
	// fmt.Println(s)
	act := &Activity{settings: s} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	settings *Settings
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	//Reading inputs
	bucket := ctx.GetInput("bucket").(string)
	item := ctx.GetInput("item").(string) // FYI - Item is the input, key is the output
	f2check := ctx.GetInput("file2Check").(string)
	region := ctx.GetInput("region").(string) // "us-east-1"
	checks3 := ctx.GetInput("checkS3").(string)
	checklocal := ctx.GetInput("checkLocal").(string)

	// defining activity wide variables (str to be returned, do we download a file, mod time of file, etc)
	var rtnfilestr string
	download := false
	var modifiedtime time.Time
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	filein := f2check
	fileout := f2check

	//Checking is local is looking at newsest in a dir or just at one file
	if checklocal == "dir" {
		ctx.Logger().Infof("checklocal is '%s' so we are checking the newest file in %s against s3", checklocal, f2check)
	} else {
		ctx.Logger().Infof("checklocal is '%s' so we are checking %s against s3", checklocal, f2check)
	}

	//Checking is local is looking at newsest in with a prefix or just at one item
	if checks3 == "prefix" {
		ctx.Logger().Infof("checks3 is '%s' so we are comparing the newest item with the prefix %s against local", checks3, item)
	} else {
		ctx.Logger().Infof("checks3 is '%s' so we are comparing %s against local", checks3, item)
	}

	// looking at local dir for newest file, if no files in directory defining an in/out file to fit with the file check below
	if checklocal == "dir" {
		// fileout=file+
		files, err := ioutil.ReadDir((filein))
		if err != nil {
			return false, fmt.Errorf("Unable to read files in directory %s", f2check)
		}
		for _, f := range files {
			t := f.ModTime()
			if t.After(modifiedtime) {
				modifiedtime = t
				filein = filein + "/" + f.Name()
				fileout = fileout + "/" + timestamp
			}
		}
		if len(files) == 0 {
			ctx.Logger().Infof("no file in %s", filein)
			fileout = fileout + "/" + timestamp
			filein = fileout
		}
	}

	// checking if file exists and then chenging download at modtime appropriately
	if file, err := os.Stat(filein); os.IsNotExist(err) {
		ctx.Logger().Infof("%s does not exist, must download new file", filein)
		download = true
	} else {
		loc, _ := time.LoadLocation("UTC")
		modifiedtime = file.ModTime().In(loc)
		ctx.Logger().Infof("%s was last modified at %s", filein, modifiedtime)
	}

	//setting up s3 env
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return true, fmt.Errorf("Unable to start aws session, %v", err)
	}
	svc := s3.New(sess)

	// checking if we are looking at s3 prefix or s3 item and checking times
	// FYI - Item is the input, key is the output
	var key string
	if checks3 == "item" {

		input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		}

		s3item, err := svc.GetObject(input)
		if err != nil {
			return false, fmt.Errorf("Unable to get info on item %s: %s", item, err)
		}

		t := *s3item.LastModified
		if t.After(modifiedtime) {
			key = item
			modifiedtime = t
			download = true
		}

	} else if checks3 == "prefix" {

		params := &s3.ListObjectsV2Input{
			Bucket: aws.String("flogo-ml"),
			Prefix: aws.String(item),
		}

		resp, err := svc.ListObjectsV2(params)
		if err != nil {
			return true, fmt.Errorf("Unable to list items in bucket %q, %v", bucket, err)
		}

		for _, listitem := range resp.Contents {
			t := *listitem.LastModified
			k := *listitem.Key
			if t.After(modifiedtime) && k[len(k)-1:] != "/" {
				key = k
				modifiedtime = t
				download = true
			}
		}
	}

	// downloading file or not
	if download {
		ctx.Logger().Infof("downloading %s modified at %s", key, modifiedtime)

		file, err := os.Create(fileout)

		if err != nil {
			return false, fmt.Errorf("Unable to download item %q, %v", item, err)
		}

		downloader := s3manager.NewDownloader(sess)
		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
		if err != nil {
			return false, fmt.Errorf("Unable to download item %q, %v", item, err)
		}

		ctx.Logger().Infof("Sucessfully downladed %s as %s", key, fileout)
		rtnfilestr = fileout
	} else {
		ctx.Logger().Infof("Newer file not found, no need to download file.")
		// fileout = f2check
		if checklocal == "file" {
			rtnfilestr = filein
		} else {
			rtnfilestr = "noFile"
		}

	}

	// returning string of file that is newest (either downloaded or if it is already newest)
	output := &Output{ModelFile: rtnfilestr}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
