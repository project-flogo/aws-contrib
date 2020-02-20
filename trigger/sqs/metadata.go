package sqs

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Region string `md:"region"`
}

type HandlerSettings struct {
	QueueURL		string	`md:"queueUrl"`
	WaitTimeSeconds int64 	`md:"waitTime"`
}

type Output struct {
	Data string `md:"data"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data":o.Data
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error

	o.Data , err = coerce.ToString(values["data"])
	
	return err
}
/*
func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
	
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	var err error
	

	return nil
}
*/