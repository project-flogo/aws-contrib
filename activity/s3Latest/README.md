# AWS s3Latest
This activity allows you to check if an s3 file is newer than a local file and if so pull it down.  The activity returns the name of the file that is newest (whether it is the original local file or the downloaded file).

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/aws-contrib/activity/s3Latest
```
## Configuration
To configure AWS credentials see [configuring-sdk](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).  Some of the possible options arefor the AWS credentials are either [env variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) or [set up the aws cli environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)

### Settings:
| Name      | Type   | Description
|:---       | :---   | :---     

### Input:
| Name      | Type   | Description
|:---       | :---   | :---     
| subject   | string | The message subject
| message   | any    | The message, either a string, object or params
| Bucket    |string  | AWS S3 bucket
| Item      |string  | AWS item/prefix to check
| File2Check|string  | local file/path to check
| Region    |string  | AWS Region
| CheckLocal|string  | either 'file' or'dir' based on what is being checked locally
| CheckS3   |string  | either 'item' or'prefix' depending on whether s3 is checking a specific item or a prefix

### Output:
| Name      | Type   | Description
|:---       | :---   | :---     
| modelFile | string | the file that is newest


## Examples
Coming soon...

