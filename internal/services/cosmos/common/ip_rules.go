package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func CosmosDBIpRulesToIpRangeFilter(ipRules *[]documentdb.IPAddressOrRange) []string {
	ipRangeFilter := make([]string, 0)
	if ipRules != nil {
		for _, ipRule := range *ipRules {
			ipRangeFilter = append(ipRangeFilter, *ipRule.IPAddressOrRange)
		}
	}

	return ipRangeFilter
}

func CosmosDBIpRangeFilterToIpRules(ipRangeFilter []string) *[]documentdb.IPAddressOrRange {
	ipRules := make([]documentdb.IPAddressOrRange, 0)
	for _, ipRange := range ipRangeFilter {
		ipRules = append(ipRules, documentdb.IPAddressOrRange{
			IPAddressOrRange: utils.String(ipRange),
		})
	}

	return &ipRules
}
