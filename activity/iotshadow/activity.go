package iotshadow

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"strings"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is an Activity that is used to update an Aws IoT device shadow
// settings : {thingName, op}
// input    : {desired,reported}
// output   : {result}
type Activity struct {
	settings *Settings
	client   *iotdataplane.IoTDataPlane
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

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

		act.client = iotdataplane.New(sess, aws.NewConfig().WithRegion(region))
	} else {
		act.client = iotdataplane.New(sess)
	}

	return act, nil
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Invokes a Aws Iot Shadow Update
func (a *Activity) Eval(context activity.Context) (done bool, err error) {

	var payload []byte

	switch a.settings.Op {
	case "update":

		in := &Input{}
		err := context.GetInputObject(in)
		if err != nil {
			return false, err
		}

		req := &ShadowRequest{State: &ShadowState{}}

		if in.Desired != nil {
			req.State.Desired = in.Desired
		}

		if in.Reported != nil {
			req.State.Reported = in.Reported
		}

		reqJSON, err := json.Marshal(req)

		sInput := &iotdataplane.UpdateThingShadowInput{}
		sInput.SetThingName(a.settings.ThingName)
		sInput.SetPayload(reqJSON)
		out, err := a.client.UpdateThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	case "get":
		sInput := &iotdataplane.GetThingShadowInput{}
		sInput.SetThingName(a.settings.ThingName)
		out, err := a.client.GetThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	case "delete":

		sInput := &iotdataplane.DeleteThingShadowInput{}
		sInput.SetThingName(a.settings.ThingName)
		out, err := a.client.DeleteThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	}

	out := &Output{}
	err = json.Unmarshal(payload, &out.Result)
	if err != nil {
		return false, err
	}

	err = context.SetOutputObject(out)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ShadowRequest is a simple structure representing a Aws Shadow Update Request
type ShadowRequest struct {
	State *ShadowState `json:"state"`
}

// ShadowState is the state to be updated
type ShadowState struct {
	Desired  map[string]interface{} `json:"desired,omitempty"`
	Reported map[string]interface{} `json:"reported,omitempty"`
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
