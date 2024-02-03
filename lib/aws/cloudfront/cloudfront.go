package dns

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/stretchr/testify/require"
)

func GetAWSApiCloudfrontClient(t *testing.T) *cloudfront.Client {
	awsApiCloudfrontClient := cloudfront.NewFromConfig(GetAWSConfig(t))
	return awsApiCloudfrontClient
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}

func GetDistribution(t *testing.T, distribution_id string) *cloudfront.GetDistributionOutput {
	cfDistribution, err := GetAWSApiCloudfrontClient(t).GetDistribution(context.Background(), &cloudfront.GetDistributionInput{
		Id: &distribution_id,
	})
	require.NoError(t, err)
	return cfDistribution
}
