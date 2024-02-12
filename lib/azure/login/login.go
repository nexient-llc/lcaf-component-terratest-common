package login

import (
	"os"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func GetEnvironmentVariables() map[string]string {
	return map[string]string{
		"clientID":       os.Getenv("CLIENT_ID"),
		"clientSecret":   os.Getenv("CLIENT_SECRET"),
		"tenantID":       os.Getenv("TENANT_ID"),
		"subscriptionID": os.Getenv("SUBSCRIPTION_ID"),
	}
}

func GetServicePrincipalToken(clientID, clientSecret, tenantID string) (*adal.ServicePrincipalToken, error) {
	config := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
	return config.ServicePrincipalToken()
}
