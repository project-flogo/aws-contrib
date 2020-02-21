package sqs

import (
	"encoding/json"
	"testing"
	"sync"

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
				"queueUrl": "https://sqs.us-east-1.amazonaws.com/011182393636/sample",
				"waitTime": 5.0

				
			}
	  }
	]
	
  }`

func TestSQS_Trigger(t *testing.T) {
	f := &Factory{}
	var wg sync.WaitGroup

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}
	wg.Add(1)
 
	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)
	
	go trg.Start()
	wg.Wait()
	
}