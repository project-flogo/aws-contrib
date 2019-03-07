package lambda

import "github.com/project-flogo/core/data/coerce"

type Output struct {
	Context interface{} `md:"context"`
	Event   interface{} `md:"evt"`
}

type Reply struct {
	Data   interface{} `md:"data"`
	Status int         `md:"status"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"context": o.Context,
		"evt":     o.Event,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Context = values["context"]

	o.Event = values["evt"]

	return nil
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
