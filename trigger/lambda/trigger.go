package lambda

import (
	"context"
	"encoding/json"
	"flag"
	syslog "log"

	"github.com/project-flogo/core/support/log"

	"github.com/project-flogo/core/trigger"

	// Import the aws-lambda-go. Required for dep to pull on app create
	_ "github.com/aws/aws-lambda-go/lambda"
)

var triggerMd = trigger.NewMetadata(&Output{}, &Reply{})
var singleton *LambdaTrigger

func init() {
	trigger.Register(&LambdaTrigger{}, &LambdaFactory{})
}

// LambdaFactory AWS Lambda Trigger factory
type LambdaFactory struct {
}

//New Creates a new trigger instance for a given id
func (t *LambdaFactory) New(config *trigger.Config) (trigger.Trigger, error) {
	singleton = &LambdaTrigger{}
	return singleton, nil

}

// Metadata implements trigger.Trigger.Metadata
func (t *LambdaFactory) Metadata() *trigger.Metadata {
	return triggerMd
}

// LambdaTrigger AWS Lambda trigger struct
type LambdaTrigger struct {
	id       string
	log      log.Logger
	handlers []trigger.Handler
}

func (t *LambdaTrigger) Initialize(ctx trigger.InitContext) error {
	t.id = "Lambda"
	t.log = ctx.Logger()
	t.handlers = ctx.GetHandlers()
	return nil
}

// Invoke starts the trigger and invokes the action registered in the handler
func Invoke() (map[string]interface{}, error) {

	log.RootLogger().Info("Starting AWS Lambda Trigger")
	syslog.Println("Starting AWS Lambda Trigger..")

	// Parse the flags
	flag.Parse()

	// Looking up the arguments
	evtArg := flag.Lookup("evt")
	var evt Event
	// Unmarshall evt
	if err := json.Unmarshal([]byte(evtArg.Value.String()), &evt); err != nil {
		return nil, err
	}

	log.RootLogger().Debugf("Received evt: '%+v'\n", evt)
	syslog.Printf("Received evt: '%+v'\n", evt)

	// Get the context
	ctxArg := flag.Lookup("ctx")
	var lambdaCtx interface{}

	// Unmarshal ctx
	if err := json.Unmarshal([]byte(ctxArg.Value.String()), &lambdaCtx); err != nil {
		return nil, err
	}

	log.RootLogger().Debugf("Received ctx: '%+v'\n", lambdaCtx)
	syslog.Printf("Received ctx: '%+v'\n", lambdaCtx)
	syslog.Printf("Singleton '%+v' \n", singleton)
	log.RootLogger().Info("handlers ", singleton.handlers[0])
	//select handler, use 0th for now
	handler := singleton.handlers[0]

	out := &Output{}

	out.Context = lambdaCtx
	out.Event = evt

	results, err := handler.Handle(context.Background(), out)

	if err != nil {
		log.RootLogger().Debugf("Lambda Trigger Error: %s", err.Error())
		return nil, err
	}

	reply := Reply{}

	reply.FromMap(results)

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

type Event struct {
	Payload interface{}     `json:"payload"`
	Flogo   json.RawMessage `json:"flogo"`
}
