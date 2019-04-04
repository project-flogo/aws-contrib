package main

import (
	"context"
	"encoding/json"
	syslog "log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fl "github.com/project-flogo/aws-contrib/trigger/lambda"
)

type LambdaHandler struct {

}

func (*LambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error){

	details, err := fl.ExtractRequestDetails(ctx, payload)
	if err != nil {
		return nil, err
	}

	result, err := fl.Invoke(details)
	if err != nil {
		return nil, err
	}

	return coerceResponse(details, result)
}

func coerceResponse(details *fl.RequestDetails, result map[string]interface{}) ([]byte, error) {

	responseData := result["data"]
	statusCode := result["status"].(int)

	var responseRaw []byte

	if val, ok := responseData.(string); ok {
		responseRaw = []byte(val)
	} else {
		var err error
		responseRaw, err = json.Marshal(responseData)
		if err != nil {
			return nil, err
		}
	}

	var resultBytes []byte

	// Check if API GW request. If so, build the correct response
	switch details.EventType {
	case fl.EtAwsApiGw:
		resp := events.APIGatewayProxyResponse{
			StatusCode: func() int {
				if statusCode == 0 {
					return 200
				} else {
					return statusCode
				}
			}(),
			Body: func() string {
				return string(responseRaw)
			}(),
			IsBase64Encoded: false,
		}
		var err error
		resultBytes, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}

	default:
		resultBytes = responseRaw
	}

	return resultBytes, nil
}

func main() {
	syslog.Println("Starting AWS Lambda Trigger..")
	lambda.Start(&LambdaHandler{})
}


