package s3newestmodel

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {
	// settings := &Settings{ReplaceFile: true, UnZip: false}
	// iCtx := test.NewActivityInitContext(settings, nil)
	// act, err := New(iCtx)

	// tc := test.NewActivityContext(act.Metadata())
	// var p []interface{}
	// p = append(p, 0)
	// input := &Input{
	// 	Bucket: "flogo-ml",
	// 	// Item:       "model_tests/Archive_20190315.zip",
	// 	File2Check: "key.zip",
	// 	CheckLocal: "file",
	// 	// CheckS3:    "item",
	// 	Item: "model_tests/",
	// 	// File2Check: "model_dir",
	// 	// CheckLocal: "dir",
	// 	CheckS3: "prefix",

	// 	Region: "us-east-1",
	// }
	// err = tc.SetInputObject(input)
	// assert.Nil(t, err)

	// done, err := act.Eval(tc)
	// assert.True(t, done)
	// assert.Nil(t, err)

	// output := &Output{}
	// tc.GetOutputObject(output)
	// fmt.Println("name of newest file is returned:", output.ModelFile)
	// // assert.Nil(t, err)
	// // assert.Equal(t, "data has been inserted into database", output.Output)
}
