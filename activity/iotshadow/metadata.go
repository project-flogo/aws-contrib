package iotshadow

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ThingName string `md:"thingName,required"`                      // The name of the "thing" in Aws IoT
	Op        string `md:"op,required,allowed(get,update,delete)"`  // The Aws IoT shadow operation to perform
}

type Input struct {
	Desired  map[string]interface{} `md:"desired"`  // The desired state of the thing
	Reported map[string]interface{} `md:"reported"` // The reported state of the thing
}

type Output struct {
	Result   map[string]interface{} `md:"result"`  // The response shadow document
}

func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"desired":  o.Desired,
		"reported": o.Reported,
	}
}

func (o *Input) FromMap(values map[string]interface{}) error {

	var err error
	o.Desired, err = coerce.ToObject(values["desired"])
	if err != nil {
		return err
	}
	o.Reported, err = coerce.ToObject(values["reported"])
	if err != nil {
		return err
	}
	return nil
}



func (r *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": r.Result,
	}
}

func (r *Output) FromMap(values map[string]interface{}) error {

	var err error
	r.Result, err = coerce.ToObject(values["result"])
	if err != nil {
		return err
	}

	return nil
}
