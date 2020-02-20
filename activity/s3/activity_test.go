package s3

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}
func TestUploadAction(t *testing.T) {
	settings := &Settings{AWSRegion: "us-east-2"}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	input := &Input{Action: "upload", S3BucketName: "", S3Location: "foogle", LocalLocation: "activity.go" }

	tc.SetInputObject(input)

	_, err = act.Eval(tc)

	assert.Nil(t,err)
}
func TestDownloadAction(t *testing.T) {
	settings := &Settings{AWSRegion: "us-east-2"}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	input := &Input{Action: "download", S3BucketName: "", S3Location: "foogle", LocalLocation: "" }

	tc.SetInputObject(input)

	_, err = act.Eval(tc)

	assert.Nil(t,err)
}

func TestDeleteAction(t *testing.T) {
	settings := &Settings{AWSRegion: "us-east-2"}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	input := &Input{Action: "delete", S3BucketName: "", S3Location: "foogle" }

	tc.SetInputObject(input)

	_, err = act.Eval(tc)

	assert.Nil(t,err)
}