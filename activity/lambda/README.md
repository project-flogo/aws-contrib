# Trigger Lambda function
This activity allows you to invoke an AWS Lambda function.

## Installation
### Flogo CLI
```bash
flogo install github.com/project-flogo/aws-contrib/activity/lambda
```
## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### Settings:
| Name          | Type   | Description
|:---           | :---   | :---     
| function      | string | The name or ARN of the Lambda function - **REQUIRED**
| clientContext | object | Information about the client to pass to the function via the context
| async         | bool   | Perform async invocation
| executionLog  | bool   | Include the execution log in the response
| region        | string | The AWS region, uses environment setting by default

### Input:
| Name     | Type   | Description
|:---      | :---   | :---     
| payload  | object | The payload object

### Output:
| Name   | Type   | Description
|:---    | :---   | :---     
| status | int    | The HTTP status code
| result | object | The response from the function


## Examples
Coming soon...
