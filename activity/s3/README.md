# Amazon S3

Upload or Download files from Amazon Simple Storage Service (S3). The Key Id and Secret Access Key needs to be set as Enviornment variable.



## Installation

```bash
flogo install github.com/project-flogo/aws-contrib/activity/s3
```

Link for flogo web:

```bash
https://github.comproject-flogo/aws-contrib/activity/s3
```
## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

## Settings

| Setting            | Description    |
|:-------------------|:---------------|
| awsRegion          | The AWS region your S3 bucket is in |

## Inputs

| Input              | Description    |
|:-------------------|:---------------|
| action             | The action you want to take, either `download`, `upload`, `delete`, or `copy` |
| s3BucketName       | The name of your S3 bucket |
| s3Location         | The file location on S3, this should be a full path (like `/bla/temp.txt`) |
| localLocation      | The `localLocation` is the full path to a file (like `/bla/temp.txt`) when uploading a file or the full path to a directory (like `./tmp`) when downloading a file |
| s3NewLocation      | The new file location on S3 of you want to copy a file, this should be a full path (like `/bla/temp.txt`) |

## Ouputs

| Output    | Description    |
|:----------|:---------------|
| result    | The result will contain OK if the action was carried out successfully or will contain an error message |
