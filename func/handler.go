package function

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ajaxe/openfaas-generate-datakey/func/config"
)

// AwsCachedConfig : Cached AWS config
var AwsCachedConfig *config.AwsConfigSettings
// AwsCachedCreds : Cached credentials
var AwsCachedCreds *credentials.Credentials

// Handle a serverless request
func Handle(req []byte) string {
	cred, awsConfig, err := config.AwsFileConfig()
	if err != nil {
		panic(fmt.Errorf("not able to read aws credentials: %v", err))
	}
	AwsCachedCreds = cred
	AwsCachedConfig = awsConfig
	return fmt.Sprintf("Successfully read creds, %s", AwsCachedConfig.KeyID)
}

func generateDataKey() (string, error) {
	session, err := session.NewSession(&aws.Config {
		Credentials: AwsCachedCreds,
		Region: aws.String(AwsCachedConfig.Region),
	},)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	svc := kms.New(session)
	svc.GenerateDataKeyWithoutPlaintext(&kms.GenerateDataKeyWithoutPlaintextInput {

	})
}
