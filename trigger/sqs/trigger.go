package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&HandlerSettings{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

func (t *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	var awsSQS *sqs.SQS

	if s.Region != "" {
		region, err := getRegion(s.Region)
		if err != nil {
			return nil, err
		}
		awsSQS = sqs.New(sess, aws.NewConfig().WithRegion(region))
	} else {
		awsSQS = sqs.New(sess)
	}

	return &Trigger{awsSQS: awsSQS}, nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

type Trigger struct {
	awsSQS    *sqs.SQS
	handlers   []Handler
	
}
type Handler struct {
	handler trigger.Handler
	sqsInput *sqs.ReceiveMessageInput
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	
	for _, handler := range ctx.GetHandlers() {
		sqsInput := &sqs.ReceiveMessageInput{}
		s := &HandlerSettings{}
		
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		
		t.handlers = append(t.handlers, Handler{sqsInput : sqsInput.SetQueueUrl(s.QueueURL),handler: handler})
	
	}

	return nil

}

func (h *Handler) subscribe(sqs *sqs.SQS) {
	
	output := &Output{}
	out, err := sqs.ReceiveMessage(h.sqsInput)

	output.Data = out.String()

	_, err = h.handler.Handle(context.Background(), out)

	if err != nil{
		logger.Debugf("Error while executing action %v", err.Error())
	}
	
}

func (t *Trigger) Start() error {
	
	for _, handler := range t.handlers {
		go handler.subscribe(t.awsSQS)
	}
	
	return nil
	
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	return nil
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
