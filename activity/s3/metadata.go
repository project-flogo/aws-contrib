package s3

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	AWSRegion          string `md:"awsRegion"`
}

type Input struct {
	Action             string `md:"action"`
	S3BucketName       string `md:"s3BucketName"`
	S3Location         string `md:"s3Location"`
	LocalLocation      string `md:"localLocation"`
	S3NewLocation      string `md:"s3NewLocation"`
}

type Output struct {
	Result string `md:"result"` // The message id
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"action":             i.Action,
		"s3BucketName":       i.S3BucketName,
		"s3Location":         i.S3Location,
		"localLocation":      i.LocalLocation,
		"s3NewLocation":      i.S3NewLocation,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Action, err = coerce.ToString(values["action"])
	if err != nil {
		return err
	}
	i.S3BucketName, err = coerce.ToString(values["s3BucketName"])
	if err != nil {
		return err
	}
	i.S3Location, err = coerce.ToString(values["s3Location"])
	if err != nil {
		return err
	}
	i.LocalLocation, err = coerce.ToString(values["localLocation"])
	if err != nil {
		return err
	}
	i.S3NewLocation, err = coerce.ToString(values["s3NewLocation"])
	if err != nil {
		return err
	}
	return err
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Result, err = coerce.ToString(values["result"])
	return err
}
