package dns

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

func GetDNSZonesClient(spt *adal.ServicePrincipalToken, subscriptionID string) dns.ZonesClient {
	dnsZonesClient := dns.NewZonesClient(subscriptionID)
	dnsZonesClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return dnsZonesClient
}
