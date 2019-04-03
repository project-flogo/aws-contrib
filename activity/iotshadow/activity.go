package iotshadow

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is an Activity that is used to update an Aws IoT device shadow
// settings : {thingName, op}
// input    : {desired,reported}
// output   : {result}
type Activity struct {
	settings  *Settings
	dataPlane *iotdataplane.IoTDataPlane
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

	idp := iotdataplane.New(sess)

	act := &Activity{settings: s, dataPlane: idp}

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
		out, err := a.dataPlane.UpdateThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	case "get":
		sInput := &iotdataplane.GetThingShadowInput{}
		sInput.SetThingName(a.settings.ThingName)
		out, err := a.dataPlane.GetThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	case "delete":

		sInput := &iotdataplane.DeleteThingShadowInput{}
		sInput.SetThingName(a.settings.ThingName)
		out, err := a.dataPlane.DeleteThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	}

	out := &Output{}
	//var result interface{}
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
