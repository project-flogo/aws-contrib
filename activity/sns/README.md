
# Amazon SNS
This activity allows you to send messages using Amazon SNS.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/aws-contrib/activity/sns
```
## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### Settings:
| Name     | Type   | Description
|:---      | :---   | :---     
| topicARN | string | The topic ARN - **REQUIRED**
| region   | string | The AWS region, uses environment setting by default
| json     | bool   | Use json message structure

### Input:
| Name     | Type   | Description
|:---      | :---   | :---     
| subject  | string | The message subject
| message  | any    | The message, either a string, object or params

### Output:
| Name      | Type   | Description
|:---       | :---   | :---     
| messageId | string | The message id


## Message
The message it typically a string.  Any value passed to message will be converted to a string.

### Json Message Structure
When you enable *json*, then you can provide messages for different protocols. Note that the
values for each provided protocol will be converted to a string.

| Protocol     | Description
|:---       | :---  
| default      | The default message is required, if not provided a dummy message will be created 
| email        | A message for email 
| email - json | A Message for email in JSON format
| http         | A message for HTTP 
| https        | A message for HTTPS 
| sqs          | A message for Amazon SQS"


## Examples
Coming soon...
