package function

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// AwsConfigSettings : AWS credential config
type AwsConfigSettings struct {
	KeyID, Secret, Region, KmsKeyID string
}

// AwsCachedConfig : Cached AWS config
var AwsCachedConfig *AwsConfigSettings

// AwsCachedCreds : Cached credentials
var AwsCachedCreds *credentials.Credentials

// Handle a serverless request
func Handle(req []byte) string {
	cred, awsConfig, err := AwsFileConfig()
	if err != nil {
		panic(fmt.Errorf("not able to read aws credentials: %v", err))
	}
	AwsCachedCreds = cred
	AwsCachedConfig = awsConfig
	dataKey, err := generateDataKey()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	return dataKey
}

func generateDataKey() (string, error) {
	session, err := session.NewSession(&aws.Config{
		Credentials: AwsCachedCreds,
		Region:      aws.String(AwsCachedConfig.Region),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	svc := kms.New(session)
	output, err := svc.GenerateDataKeyWithoutPlaintext(&kms.GenerateDataKeyWithoutPlaintextInput{
		KeyId:         aws.String(AwsCachedConfig.KmsKeyID),
		NumberOfBytes: aws.Int64(32),
	})
	if err != nil {
		return "", fmt.Errorf("error genreating datakey: %v", err)
	}
	return base64.StdEncoding.EncodeToString(output.CiphertextBlob), nil
}

// AwsFileConfig : reads aws config from "./aws_credentials" file
func AwsFileConfig() (*credentials.Credentials, *AwsConfigSettings, error) {
	const credFilePath = "./secrets/aws_config.json"
	data, err := ioutil.ReadFile(credFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read data from file [%s]: %v", credFilePath, err)
	}
	cfg := new(AwsConfigSettings)
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, nil, fmt.Errorf("error parsing json data")
	}
	return credentials.NewStaticCredentials(cfg.KeyID, cfg.Secret, ""), cfg, nil
}
