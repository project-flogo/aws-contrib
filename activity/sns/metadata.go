package sns

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	TopicARN string `md:"topic,required"` // The topic ARN
	Json     bool   `md:"json"`           // Use json message structure
	Region   string `md:"region"`         // The AWS region, uses environment setting by default
}

type Input struct {
	Subject string      `md:"subject"` // The message subject
	Message interface{} `md:"message"` // The message, either a string, object or params
}

type Output struct {
	MessageId string `md:"messageId"` // The message id
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"subject": i.Subject,
		"message": i.Message,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Subject, err = coerce.ToString(values["subject"])
	if err != nil {
		return err
	}
	i.Message = values["message"]
	return err
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"messageId": o.MessageId,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.MessageId, err = coerce.ToString(values["messageId"])
	return err
}
