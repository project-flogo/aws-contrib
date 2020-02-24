# AWS SQS Trigger

The SQS trigger provides your Flogo application the ability to read data from SQ. 

## Installation

### Flogo CLI

```bash
flogo install github.com/project-flogo/aws-contrib/trigger/sqs
```

## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### Settings:
| Name      | Type   | Description
|:---       | :---   | :---          
| region    | string | The region of SQS trigger.

### Handler Settings:
| Name      | Type   | Description
|:---       | :---   | :---          
| queueUrl  | string | The url of the SQS queue

### Output:
| Name      | Type   | Description
|:---       | :---   | :---        
| data      | array  | The array containing messages from SQS.
