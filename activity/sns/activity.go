package sns

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

const (
	ovMessageId = "messageId"
)

// Activity is an activity that is used to invoke a lambda function
type Activity struct {
	settings *Settings
	client   *sns.SNS
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
		act.client = sns.New(sess, aws.NewConfig().WithRegion(s.Region))
	} else {
		act.client = sns.New(sess)
	}

	return act, nil
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	ctx.Logger().Debugf("Sending SNS Message To: %s", a.settings.TopicARN)

	in := &Input{}
	err = ctx.GetInputObject(in)
	if err != nil {
		return false, err
	}

	pInput := &sns.PublishInput{TopicArn: &a.settings.TopicARN}

	var msg string
	if a.settings.Json {
		pInput.SetMessageStructure("json")

		switch t := in.Message.(type) {
		case map[string]string:
			if _, exists := t["default"]; !exists {
				t["default"] = "default message"
			}
			msg, err = coerce.ToString(t)
		case map[string]interface{}:
			if _, exists := t["default"]; !exists {
				t["default"] = "default message"
			}
			msg, err = coerce.ToString(t)
		case string:
			msg = fmt.Sprintf("{\"default\":\"%s\"", t)
		default:
			def, err := coerce.ToString(t)
			if err != nil {
				return false, err
			}
			msg = fmt.Sprintf("{\"default\":\"%s\"", def)
		}
	} else {
		msg, err = coerce.ToString(in.Message)
	}
	if err != nil {
		return false, err
	}

	pInput.SetMessage(msg)
	if in.Subject != "" {
		pInput.SetSubject(in.Subject)
	}

	if ctx.Logger().TraceEnabled() {
		if in.Subject != "" {
			ctx.Logger().Tracef("Subject: '%s'", in.Subject)
		}
		ctx.Logger().Tracef("Message: '%s'", msg)
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
