package network

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

func GetRouteTablesClient(spt *adal.ServicePrincipalToken, subscriptionID string) network.RouteTablesClient {
	routeTableClient := network.NewRouteTablesClient(subscriptionID)
	routeTableClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return routeTableClient
}

func GetSubnetsClient(spt *adal.ServicePrincipalToken, subscriptionID string) network.SubnetsClient {
	subnetsClient := network.NewSubnetsClient(subscriptionID)
	subnetsClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return subnetsClient
}

func GetFirewallsClient(spt *adal.ServicePrincipalToken, subscriptionID string) network.AzureFirewallsClient {
	firewallsClient := network.NewAzureFirewallsClient(subscriptionID)
	firewallsClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return firewallsClient
}
func GetNsgClient(spt *adal.ServicePrincipalToken, subscriptionID string) network.SecurityGroupsClient {
	nsgClient := network.NewSecurityGroupsClient(subscriptionID)
	nsgClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return nsgClient
}
