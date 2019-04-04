package lambda

import (
	"context"
	"encoding/json"
	syslog "log"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&HandlerSettings{}, &Output{}, &Reply{})
var singleton *LambdaTrigger

func init() {
	_ = trigger.Register(&LambdaTrigger{}, &LambdaFactory{})
}

// LambdaFactory AWS Lambda Trigger factory
type LambdaFactory struct {
}

//New Creates a new trigger instance for a given id
func (t *LambdaFactory) New(config *trigger.Config) (trigger.Trigger, error) {

	if singleton == nil {
		singleton = &LambdaTrigger{}
		return singleton, nil
	}

	log.RootLogger().Warn("Only one lambda trigger instance can be instantiated")

	return nil, nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *LambdaFactory) Metadata() *trigger.Metadata {
	return triggerMd
}

// LambdaTrigger AWS Lambda trigger struct
type LambdaTrigger struct {
	id             string
	log            log.Logger
	handlers       map[string]trigger.Handler
	defaultHandler trigger.Handler
}

func (t *LambdaTrigger) Initialize(ctx trigger.InitContext) error {
	t.id = "Lambda"
	t.log = ctx.Logger()
	t.defaultHandler = ctx.GetHandlers()[0]
	t.handlers = make(map[string]trigger.Handler)

	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		if s.EventType == "" {
			if t.defaultHandler == nil {
				t.defaultHandler = handler
			} else {
				log.RootLogger().Warn("Only one default handler will be used")
			}
			continue
		}

		if _, exists := t.handlers[s.EventType]; exists {
			log.RootLogger().Warnf("Only first handler for eventType '%s' will be used", s.EventType)
			break
		}

		t.handlers[s.EventType] = handler
	}

	return nil
}

// Invoke starts the trigger and invokes the action registered in the handler
func Invoke(details *RequestDetails) (map[string]interface{}, error) {

	syslog.Printf("Received request: %s\n", details.CtxInfo["awsRequestId"])

	//todo figure out how to support flogo logging in Lambda
	//log.RootLogger().Debugf("Received ctx: '%+v'\n", lambdaCtx)

	evtTypeStr := FromoEventType(details.EventType)

	syslog.Printf("Payload Type: %s\n", evtTypeStr)
	syslog.Printf("Payload: '%+v'\n", details.Event)

	out := &Output{}
	out.Context = details.CtxInfo
	out.Event = details.Event

	if details.EventType == EtFlogoOnDemand {

		// todo add event type to flogo events?
		var evt FlogoEvent
		if err := json.Unmarshal(details.Payload, &evt); err != nil {
			return nil, err
		}

		out.Event = map[string]interface{}{"payload":evt.Payload, "flogo":evt.Flogo}
	}

	out.EventType = evtTypeStr

	//select handler for the specified eventType
	handler := singleton.handlers[evtTypeStr]
	if handler == nil {
		handler = singleton.defaultHandler
	}
	results, err := handler.Handle(context.Background(), out)
	if err != nil {
		log.RootLogger().Debugf("Lambda Trigger Error: %s", err.Error())
		syslog.Printf("Lambda Trigger Error: %s", err.Error())
		return nil, err
	}

	reply := Reply{}
	err = reply.FromMap(results)
	if err != nil {
		return nil, err
	}

	if reply.Data != nil {
		if reply.Status == 0 {
			reply.Status = 200
		}
	}

	return reply.ToMap(), err
}

func (t *LambdaTrigger) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (t *LambdaTrigger) Stop() error {
	return nil
}

type FlogoEvent struct {
	Payload interface{}     `json:"payload"`
	Flogo   json.RawMessage `json:"flogo"`
}
