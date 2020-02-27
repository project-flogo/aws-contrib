package sns

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

const (
	ovMessageId = "messageId"
)
var logger log.Logger
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
		region, err := getRegion(s.Region)
		if err != nil {
			return nil, err
		}

		act.client = sns.New(sess, aws.NewConfig().WithRegion(region))
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

	logger = ctx.Logger()

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
		if reqerr, ok := err.(awserr.RequestFailure); ok {
			logger.Debug("Request failed", reqerr.Code(), reqerr.Message(), reqerr.RequestID())
		} else {
			logger.Debug("Error:", err.Error())
		}
		return false, err
	}

	err = ctx.SetOutput(ovMessageId, *pOutput.MessageId)
	if err != nil {
		return false, err
	}
	ctx.Logger().Debugf("Message sent: %s", *pOutput.MessageId)

	return true, nil
}

func getRegion(regionSetting string) (string, error) {

	var awsRegions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "cn-north-1", "cn-northwest-1", "eu-central-1", "eu-north-1", "eu-west-1", "eu-west-2", "eu-west-2", "sa-east-1"}

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
