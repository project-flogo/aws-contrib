package lambda

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Function      string                 `md:"function,required"` // The name or ARN of the Lambda function
	ClientContext map[string]interface{} `md:"clientContext"`     // Information about the client to pass to the function via the context
	Async         bool                   `md:"async"`             // Perform async invocation
	ExecutionLog  bool                   `md:"executionLog"`      // Include the execution log in the response
	Region        string                 `md:"region"`            // The AWS region, used environment setting by default
}

type Input struct {
	Payload map[string]interface{} `md:"payload"` // The payload object
}

type Output struct {
	Result map[string]interface{} `md:"result"` // The response from the function
	Status int                    `md:"status"` // The HTTP status code
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"payload": i.Payload,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Payload, err = coerce.ToObject(values["payload"])
	return err
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
		"status": o.Status,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Result, err = coerce.ToObject(values["result"])
	if err != nil {
		return err
	}

	o.Status, err = coerce.ToInt(values["status"])

	return nil
}
