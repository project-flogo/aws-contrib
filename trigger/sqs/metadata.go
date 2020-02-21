package sqs

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Region string `md:"region"`
}

type HandlerSettings struct {
	QueueURL		string	`md:"queueUrl,required"`
}

type Output struct {
	Data []interface{} `md:"data"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data":o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error

	o.Data , err = coerce.ToArray(values["data"])
	
	return err
}
