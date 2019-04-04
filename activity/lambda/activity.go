package lambda

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is an activity that is used to invoke a lambda function
type Activity struct {
	settings      *Settings
	client        *lambda.Lambda
	clientContext string
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
		act.client = lambda.New(sess, aws.NewConfig().WithRegion(s.Region))
	} else {
		act.client = lambda.New(sess)
	}

	if s.ClientContext != nil {
		var b []byte
		b, err := json.Marshal(&s.ClientContext)
		if err != nil {
			return nil, err
		}

		base64.StdEncoding.EncodeToString(b)
		act.clientContext = string(b)
	}

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

	iInput := &lambda.InvokeInput{FunctionName: &a.settings.Function}

	if a.clientContext != "" {
		iInput.SetClientContext(a.clientContext)
	}

	if a.settings.Async {
		iInput.SetInvocationType(lambda.InvocationTypeEvent)
	}

	if a.settings.ExecutionLog {
		iInput.SetLogType(lambda.LogTypeTail)
	}

	if in.Payload != nil {
		b, err := json.Marshal(&in.Payload)
		if err != nil {
			return false, err
		}
		iInput.SetPayload(b)
	}

	iOutput, err := a.client.Invoke(iInput)
	if err != nil {
		ctx.Logger().Tracef("Lambda invoke error: %v", err)
		return false, err
	}
	ctx.Logger().Tracef("Lambda response: %s", string(iOutput.Payload))

	out := &Output{}

	err = json.Unmarshal(iOutput.Payload, &out.Result)
	if err != nil {
		return false, err
	}

	out.Status = int(*iOutput.StatusCode)

	err = ctx.SetOutputObject(out)
	if err != nil {
		return false, err
	}

	return true, nil
}
