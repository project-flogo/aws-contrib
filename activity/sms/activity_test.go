package sms

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

func TestSimpleSMS(t *testing.T) {

	settings := &Settings{Region: "us-east-1", SmsType: "Promotional", SenderID: "", MaxPrice:0.0 }

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	tc.SetInput("to", "+99999")
	tc.SetInput("message", "hello world")

	_, err = act.Eval(tc)
	assert.Nil(t, err)
}

