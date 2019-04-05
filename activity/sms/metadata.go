package sms

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	SmsType  string `md:"smsType"`          // The type of SMS to send, defaults to Promotional
	SenderID string `md:"senderID"`         // The Sender ID for the SMS (note: not supported in all countries)
	Region   string `md:"region, required"` // The AWS region to use (note: SMS is not supported in all AWS regions)
	MaxPrice float64 `md:"maxPrice"`        // The maximum amount in USD that you are willing to spend to send a message
}

type Input struct {
	To      string `md:"subject"` // The phone number (in international format) to which to send the SMS
	Message string `md:"message"` // The message to send
}

type Output struct {
	MessageId string `md:"messageId"` // The message id
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"to":      i.To,
		"message": i.Message,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.To, err = coerce.ToString(values["to"])
	if err != nil {
		return err
	}
	i.Message, err = coerce.ToString(values["message"])
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
