package lambda

import "github.com/project-flogo/core/data/coerce"

type HandlerSettings struct {
	EventType string `md:"eventType"` // The type of event to handle, (e.g. aws:s3, aws:sns)
}

type Output struct {
	Context   map[string]interface{} `md:"context"`    // The lambda context information
	Event     map[string]interface{} `md:"event"`      // The event data
	EventType string                 `md:"eventType"`  // The event type, (e.g. aws:s3, aws:sns)
}

type Reply struct {
	Data   interface{} `md:"data"`   // The data to reply with
	Status int         `md:"status"` // The status code to reply with
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"context":   o.Context,
		"event":     o.Event,
		"eventType": o.EventType,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Context, err = coerce.ToObject(values["context"])
	if err != nil {
		return err
	}
	o.Event, err = coerce.ToObject(values["event"])
	if err != nil {
		return err
	}
	o.EventType, err = coerce.ToString(values["eventType"])
	return err
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status": r.Status,
		"data":   r.Data,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	var err error
	r.Status, err = coerce.ToInt(values["status"])
	if err != nil {
		return err
	}
	r.Data, _ = values["data"]

	return nil
}
