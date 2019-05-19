# Generate Data Key using AWS KMS

A simple Open FaaS function to generate data key using the KMS arn provided as configuration.

## Configuration

Function requires a secrets config file at `./func/secrets/aws_config.json` with following json:

```json
{
    "KeyID": "aws-access-key-id",
    "Secret": "aws-access-key-secret",
    "Region": "us-east-1",
    "KmsKeyID": "arn-kms-key"
}
```