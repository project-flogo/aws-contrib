<!--
title: AWS Device Shadow
weight: 4605
-->
# AWS Device Shadow
This activity allows you to update/get/delete a device shadow on AWS.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/aws-contrib/activity/iotshadow
```
## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### Settings:
| Name       | Type   | Description
|:---        | :---   | :---     
| thingName  | string | The name of the "thing" in AWS IoT **REQUIRED**
| op         | string | The AWS IoT shadow operation to perform  (Allowed values are get, update, delete) - **REQUIRED**
| region     | string | The AWS region, uses environment setting by default

### Input:
| Name     | Type   | Description
|:---      | :---   | :---     
| desired  | object | The desired state of the thing
| reported | object | The reported state of the thing

### Output:
| Name   | Type   | Description
|:---    | :---   | :---     
| result | object | The response shadow document

## Examples

### Update Temp
Configure a task in flow to update the device shadow of 'raspberry-pi' with a reported temp of "50".

```json
{
  "id": "shadow_update",
  "name": "Update AWS Device Shadow",
  "activity": {
    "ref": "github.com/project-flogo/aws-contrib/",
    "settings": {
      "thingName": "raspberry-pi",
       "op": "update"
    },
    "input": {
      "reported": { "temp":"50" }
    }
  }
}
```
