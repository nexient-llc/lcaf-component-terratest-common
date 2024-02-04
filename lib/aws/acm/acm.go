package acm

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/stretchr/testify/require"
)

func GetAWSApiAcmClient(t *testing.T) *acm.Client {
	awsApiAcmClient := acm.NewFromConfig(GetAWSConfig(t))
	return awsApiAcmClient
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}

func ListCertificates(t *testing.T, certificate_arn string) *acm.ListCertificatesOutput {
	certificates, err := GetAWSApiAcmClient(t).ListCertificates(context.TODO(), &acm.ListCertificatesInput{})
	require.NoErrorf(t, err, "unable to get certificate, %v", err)
	return certificates
}
