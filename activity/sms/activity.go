package sms

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

const (
	ovMessageId = "messageId"

	defaultMaxPrice = 0.01
)

// Activity is an activity that is used to invoke a lambda function
type Activity struct {
	settings *Settings
	client   *sns.SNS
	msgAttrs map[string]*sns.MessageAttributeValue
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	// assumes session configuration via environment
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}

	if s.Region != "" {
		region, err := getRegion(s.Region)
		if err != nil {
			return nil, err
		}
		act.client = sns.New(sess, aws.NewConfig().WithRegion(region))
	} else {
		act.client = sns.New(sess)
	}

	maxPrice := defaultMaxPrice
	if s.MaxPrice != 0.0 {
		maxPrice = s.MaxPrice
	}
	maxPriceString := fmt.Sprintf("%0.5f", maxPrice)

	act.msgAttrs = getMessageAttrs(s.SmsType, s.SenderID, maxPriceString)

	return act, nil
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	in := &Input{}
	err = ctx.GetInputObject(in)
	if err != nil {
		return false, err
	}

	ctx.Logger().Debugf("Sending SMS Message To: %s", in.To)

	pInput := &sns.PublishInput{
		PhoneNumber:       aws.String(in.To),
		Message:           aws.String(in.Message),
		MessageAttributes: a.msgAttrs,
	}

	if ctx.Logger().TraceEnabled() {
		ctx.Logger().Tracef("Message: '%s'", in.Message)
	}

	pOutput, err := a.client.Publish(pInput)
	if err != nil {
		return false, err
	}

	err = ctx.SetOutput(ovMessageId, *pOutput.MessageId)
	if err != nil {
		return false, err
	}
	ctx.Logger().Debugf("Message sent: %s", *pOutput.MessageId)

	return true, nil
}

func getMessageAttrs(smsType, senderID, maxPrice string) map[string]*sns.MessageAttributeValue {

	messageAttrs := map[string]*sns.MessageAttributeValue{
		"AWS.SNS.SMS.SMSType": {
			DataType:    aws.String("String"),
			StringValue: aws.String(smsType),
		},
		"AWS.SNS.SMS.MaxPrice": {
			DataType:    aws.String("Number"),
			StringValue: aws.String(maxPrice),
		},
	}

	if senderID != "" {
		messageAttrs["AWS.SNS.SMS.SenderID"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(senderID),
		}
	}

	return messageAttrs
}

func getRegion(regionSetting string) (string, error) {

	var awsRegions = []string{"us-east-1", "us-west-2", "ap-northeast-1", "ap-southeast-1", "ap-southeast-2", "eu-west-1"}

	region := strings.ToLower(regionSetting)
	valid := false
	for _, aRegion := range awsRegions {
		if region == aRegion {
			valid = true
			break
		}
	}

	if !valid {
		return "", fmt.Errorf("unsupported region: %s", regionSetting)
	}

	return region, nil
}
