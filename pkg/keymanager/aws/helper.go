package aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

func getAWSSecretString(secretName, region string) (string, error) {
	// Create a Secrets Manager client
	sess, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		return "", err
	}
	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}
	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	length, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return "", errors.Wrap(err, "base64 decode error")
	}

	return string(decodedBinarySecretBytes[:length]), nil
}
