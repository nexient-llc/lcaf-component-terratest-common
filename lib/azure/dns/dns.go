package dns

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

func GetDNSZonesClient(spt *adal.ServicePrincipalToken, subscriptionID string) dns.ZonesClient {
	dnsZonesClient := dns.NewZonesClient(subscriptionID)
	dnsZonesClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return dnsZonesClient
}

func GetDNSZoneRecordSetsClient(spt *adal.ServicePrincipalToken, subscriptionID string) dns.RecordSetsClient {
	dnsZoneRecordSetsClient := dns.NewRecordSetsClient(subscriptionID)
	dnsZoneRecordSetsClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return dnsZoneRecordSetsClient
}

func GetPrivateDNSZoneRecordSetsClient(spt *adal.ServicePrincipalToken, subscriptionID string) privatedns.RecordSetsClient {
	privateDnsZoneRecordSetsClient := privatedns.NewRecordSetsClient(subscriptionID)
	privateDnsZoneRecordSetsClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	return privateDnsZoneRecordSetsClient
}
