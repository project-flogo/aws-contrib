package lambda

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"strings"
)

type EvtType int

const (
	EtUnknown EvtType = iota
	EtFlogoOnDemand
	EtAwsApiGw
	EtAwsAutoScaling
	EtAwsCloudFront
	EtAwsCloudWatch
	EtAwsCloudWatchLogs
	EtAwsCodeCommit
	EtAwsCodePipeline
	EtAwsCognito
	EtAwsConfig
	EtAwsDynamodb
	EtAwsKinesis
	EtAwsKinesisAnalytics
	EtAwsKinesisFirehose
	EtAwsS3
	EtAwsSes
	EtAwsSns
	EtAwsSqs
)

var toEventType = map[string]EvtType{
	"flogo:ondemand":        EtFlogoOnDemand,
	"aws:apigw":             EtAwsApiGw,
	"aws:autoscaling":       EtAwsAutoScaling,
	"aws:cloudfront":        EtAwsCloudFront,
	"aws:cloudwatch":        EtAwsCloudWatch,
	"aws:cloudwatch-logs":   EtAwsCloudWatchLogs,
	"aws:codecommit":        EtAwsCodeCommit,
	"aws:codepipeline":      EtAwsCodePipeline,
	"aws:cognito":           EtAwsCognito,
	"aws:config":            EtAwsConfig,
	"aws:dynamodb":          EtAwsDynamodb,
	"aws:kinesis":           EtAwsKinesis,
	"aws:kinesis-analytics": EtAwsKinesisAnalytics,
	"aws:kinesis-firehose":  EtAwsKinesisFirehose,
	"aws:s3":                EtAwsS3,
	"aws:ses":               EtAwsSes,
	"aws:sns":               EtAwsSns,
	"aws:sqs":               EtAwsSqs,
}

var fromEventType = map[EvtType]string{
	EtFlogoOnDemand:       "flogo:ondemand",
	EtAwsApiGw:            "aws:apigw",
	EtAwsAutoScaling:      "aws:autoscaling",
	EtAwsCloudFront:       "aws:cloudfront",
	EtAwsCloudWatch:       "aws:cloudwatch",
	EtAwsCloudWatchLogs:   "aws:cloudwatch-logs",
	EtAwsCodeCommit:       "aws:codecommit",
	EtAwsCodePipeline:     "aws:codepipeline",
	EtAwsCognito:          "aws:cognito",
	EtAwsConfig:           "aws:config",
	EtAwsDynamodb:         "aws:dynamodb",
	EtAwsKinesis:          "aws:kinesis",
	EtAwsKinesisAnalytics: "aws:kinesis-analytics",
	EtAwsKinesisFirehose:  "aws:kinesis-firehose",
	EtAwsS3:               "aws:s3",
	EtAwsSes:              "aws:ses",
	EtAwsSns:              "aws:sns",
	EtAwsSqs:              "aws:sqs",
}

func ToEventType(str string) EvtType {
	if et, exists := toEventType[str]; exists {
		return et
	}
	return EtUnknown
}

func FromoEventType(et EvtType) string {
	if et, exists := fromEventType[et]; exists {
		return et
	}
	return "unknown"
}

func GetEventType(payload map[string]interface{}) EvtType {

	if fg, exists := payload["flogo"]; exists {
		if fgPkg, ok := fg.(map[string]interface{}); ok {
			if _, exists := fgPkg["flow"]; exists {
				return EtFlogoOnDemand
			}
		}

		return EtUnknown
	}

	// look throw all AWS types that have 'Records'
	// aws:codecommit, aws:dynamodb, aws:kinesis, aws:s3, aws:ses, aws:sns, aws:sqs
	if records, exists := payload["Records"]; exists {
		if rArray, ok := records.([]interface{}); ok {
			if zeroRecord, ok := rArray[0].(map[string]interface{}); ok {
				if es, exists := zeroRecord["eventSource"]; exists {
					if eventSource, ok := es.(string); ok {
						return ToEventType(eventSource)
					} else {
						return EtUnknown
					}
				}
				if es, exists := zeroRecord["EventSource"]; exists {
					if eventSource, ok := es.(string); ok {
						return ToEventType(eventSource)
					} else {
						return EtUnknown
					}
				}
				if _, exists := zeroRecord["cf"]; exists {
					return EtAwsCloudFront
				}
			}
		}

		return EtUnknown
	}

	// look throw all AWS types that have 'records'
	// aws:kinesis-analytics, aws:kinesis-firehose"
	if records, exists := payload["records"]; exists {

		if arn, exists := payload["applicationArn"]; exists {
			if arnStr, ok := arn.(string); ok {
				if strings.HasPrefix(arnStr, "arn:aws:kinesisanalytics") {
					return EtAwsKinesisAnalytics
				}
			} else {
				return EtUnknown
			}
		}

		if rArray, ok := records.([]interface{}); ok {
			if zeroRecord, ok := rArray[0].(map[string]interface{}); ok {
				if _, exists := zeroRecord["kinesisRecordMetadata"]; exists {
					return EtAwsKinesisFirehose
				}
			}
		}

		return EtUnknown
	}

	//API Gateway
	if _, exists := payload["requestContext"]; exists {
		return EtAwsApiGw

		//custom auth request type-request
		// "type": "REQUEST",
		// "methodArn": "arn:aws:execute-api:...
	}

	if src, exists := payload["source"]; exists {
		if source, ok := src.(string); ok {
			if source == "aws.autoscaling" {
				return EtAwsAutoScaling
			}
			if source == "aws.events" {
				return EtAwsCloudWatch
			}
		}
	}

	if _, exists := payload["awslogs"]; exists {
		return EtAwsCloudWatchLogs
	}

	if _, exists := payload["CodePipeline.job"]; exists {
		return EtAwsCodePipeline
	}

	if et, exists := payload["eventType"]; exists {
		if eventType, ok := et.(string); ok {
			if eventType == "SyncTrigger" {
				return EtAwsCognito
			}
		}
	}

	if _, exists := payload["configRuleId"]; exists {
		return EtAwsConfig
	}

	return EtUnknown
}

func ExtractRequestDetails(ctx context.Context, payload []byte) (*RequestDetails, error) {

	var payloadObj map[string]interface{}
	err := json.Unmarshal(payload, &payloadObj)
	if err != nil {
		return nil, err
	}

	et := GetEventType(payloadObj)

	lambdaCtx, _ := lambdacontext.FromContext(ctx)
	ctxInfo := map[string]interface{}{
		"logStreamName":   lambdacontext.LogStreamName,
		"logGroupName":    lambdacontext.LogGroupName,
		"functionName":    lambdacontext.FunctionName,
		"functionVersion": lambdacontext.FunctionVersion,
		"awsRequestId":    lambdaCtx.AwsRequestID,
		"memoryLimitInMB": lambdacontext.MemoryLimitInMB,
	}

	return &RequestDetails{CtxInfo: ctxInfo, Event: payloadObj, Payload: payload, EventType: et}, nil
}

type RequestDetails struct {
	CtxInfo   map[string]interface{}
	Payload   []byte
	Event     map[string]interface{}
	EventType EvtType
}
