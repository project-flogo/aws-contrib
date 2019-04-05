package iotshadow

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ThingName string `md:"thingName,required"`                      // The name of the "thing" in AWS IoT
	Op        string `md:"op,required,allowed(get,update,delete)"`  // The AWS IoT shadow operation to perform
	Region    string `md:"region"`                                  // The AWS region, uses environment setting by default
}

type Input struct {
	Desired  map[string]interface{} `md:"desired"`  // The desired state of the thing
	Reported map[string]interface{} `md:"reported"` // The reported state of the thing
}

type Output struct {
	Result   map[string]interface{} `md:"result"`  // The response shadow document
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"desired":  i.Desired,
		"reported": i.Reported,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Desired, err = coerce.ToObject(values["desired"])
	if err != nil {
		return err
	}
	i.Reported, err = coerce.ToObject(values["reported"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Result, err = coerce.ToObject(values["result"])
	if err != nil {
		return err
	}

	return nil
}
