package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	fLambda "github.com/project-flogo/aws-contrib/trigger/lambda"
)

type LambdaHandler struct {

}

func (*LambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error){

	details, err := fLambda.ExtractRequestDetails(ctx, payload)
	if err != nil {
		return nil, err
	}

	result, err := fLambda.Invoke(details)
	if err != nil {
		return nil, err
	}

	return fLambda.GenerateResponse(details, result)
}

func main() {
	log.Println("Starting AWS Lambda Trigger..")
	lambda.StartHandler(&LambdaHandler{})
}


