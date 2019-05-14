<!--
title: AWS SNS-SMS
weight: 4605
-->
# AWS SNS-SMS
This activity allows you to send SMS messages using AWS SNS.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/aws-contrib/activity/sms
```
## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### Settings:
| Name     | Type    | Description
|:---      | :---    | :---     
| smsType  | string  | The type of SMS to send, defaults to Promotional
| region   | string  | The AWS region, uses environment setting by default
| senderID | string  | The Sender ID for the SMS (note: not supported in all countries)
| maxPrice | float64 | The maximum amount in USD that you are willing to spend to send a message


### Input:
| Name    | Type   | Description
|:---     | :---   | :---     
| to      | string | The phone number to which to send the SMS (e.g. +15555550100)
| message | string | The message to send

### Output:
| Name      | Type   | Description
|:---       | :---   | :---     
| messageId | string | The message id


#### Additional Information

You can find additional information at Amazon's [Supported Regions and Countries](https://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html) page.
There you can see what are the currently supported regions and also which countries support SenderID. 

## Examples
Coming soon...
