package sqs

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "flogo-sqs",
	"ref": "github.com/project-flogo/aws-contrib/trigger/sqs",
	"settings": {
	  "region": "us-east-1"
	},
	"handlers": [
	  {
			"action":{
				"id":"dummy"
			},
			"settings": {
				"queue": ""
				
			}
	  }
	]
	
  }`

func TestKafkaTrigger_Initialize(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	trg.Start()
}