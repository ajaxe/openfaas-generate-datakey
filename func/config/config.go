package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// AwsConfigSettings : AWS credential config
type AwsConfigSettings struct {
	KeyID, Secret, Region string
}

// AwsFileConfig : reads aws config from "./aws_credentials" file
func AwsFileConfig() (*credentials.Credentials, *AwsConfigSettings, error) {
	const credFilePath = "./secrets/aws_config"
	file, err := os.Open(credFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open file [%s]: %v", credFilePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	cfg := new(AwsConfigSettings)
	ctr := 0
	for scanner.Scan() {
		ctr++
		line := scanner.Text()
		if len(line) <= 0 {
			continue
		}
		kvp := strings.Split(strings.TrimSpace(line), "=")
		if len(kvp) != 2 {
			return nil, nil, fmt.Errorf("invalid config format at line %d", ctr)
		}
		key, val := strings.TrimSpace(kvp[0]), strings.TrimSpace(kvp[1])
		switch key {
		case "aws_access_key_id":
			cfg.KeyID = val
		case "aws_secret_access_key":
			cfg.Secret = val
		case "aws_Region":
			cfg.Region = val
		default:
			return nil, nil, fmt.Errorf("invalid config: %s", key)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("reading file input: %v", err)
	}
	return credentials.NewStaticCredentials(cfg.KeyID, cfg.Secret, ""), cfg, nil
}
