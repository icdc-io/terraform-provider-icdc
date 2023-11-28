package icdc

import (
	"fmt"
	"os"
	"strings"

	"github.com/3th1nk/cidr"
)

type NetworkCollection struct {
	Resources []Network `json:"resources"`
}

type SubnetCollection struct {
	Resources []Subnet `json:"resources"`
}

type Network struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	EmsRef  string   `json:"ems_ref,omitempty"`
	Subnets []Subnet `json:"cloud_subnets"`
}

type Subnet struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Cidr            string   `json:"cidr"`
	Gateway         string   `json:"gateway"`
	DnsNameservers  []string `json:"dns_nameservers"`
	EmsRef          string   `json:"ems_ref"`
	EmsId           string   `json:"ems_id"`
	NetworkProtocol string   `json:"network_protocol"`
	NetworkId       string   `json:"cloud_network_id"`
	RouterId        string   `json:"network_router_id"`
}

type addNetworkBody struct {
	Action string        `json:"action"`
	Name   string        `json:"name"`
	Mtu    int           `json:"mtu"`
	Subnet addSubnetBody `json:"subnet"`
}

type addSubnetBody struct {
	Cidr            string   `json:"cidr"`
	IpVersion       int      `json:"ip_version"`
	NetworkProtocol string   `json:"network_protocol"`
	DnsNameservers  []string `json:"dns_nameservers"`
	Name            string   `json:"name"`
}

type EmsProvider struct {
	Resources []struct {
		Id string `json:"id"`
	} `json:"resources"`
}

func convertName(name string) string {
	prefix := fmt.Sprintf("%s_%s_", os.Getenv("LOCATION"), os.Getenv("ACCOUNT"))
	if strings.HasPrefix(name, prefix) {
		name, _ = strings.CutPrefix(name, prefix)
		return name
	}
	return name
}

func fetchNetworkProtocol(cidr_pure string) (int, string) {
	c, _ := cidr.Parse(cidr_pure)

	if c.IsIPv4() {
		return 4, "ipv4"
	} else {
		return 6, "ipv6"
	}
}

func getProviderId() (string, error) {
	var emsProvider *EmsProvider
	responseBody, err := requestApi("GET", "api/compute/v1/providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager", nil)

	if err != nil {
		return "", err
	}

	err = responseBody.Decode(&emsProvider)

	if err != nil {
		return "", err
	}

	return emsProvider.Resources[0].Id, nil
}

func getNetworkObject(n string, c string, pId string) (Network, error) {
	requestUrl := fmt.Sprintf("api/compute/v1/providers/%s/cloud_networks?expand=resources&attributes=cloud_subnets", pId)
	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return Network{}, fmt.Errorf("can't fetch cloud_networks list: %s", err)
	}

	var networkCollection *NetworkCollection
	err = responseBody.Decode(&networkCollection)

	if err != nil {
		return Network{}, fmt.Errorf("can't decode cloud_networks list response: %s", err)
	}

	for _, network := range networkCollection.Resources {
		if network.Subnets[0].Cidr == c && n == convertName(network.Name) {
			return network, nil
		}
	}

	return Network{}, fmt.Errorf("network %s [%s] was not created", n, c)
}

func subnetCreated(n string, c string, pId string) (bool, error) {

	requestUrl := fmt.Sprintf("api/compute/v1/providers/%s/cloud_subnets?expand=resources", pId)

	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return false, fmt.Errorf("can't fetch cloud_networks list: %s", err)
	}
	var subnetCollection *SubnetCollection
	err = responseBody.Decode(&subnetCollection)

	if err != nil {
		return false, fmt.Errorf("can't decode cloud subnets list: %s", err)
	}

	for _, subnet := range subnetCollection.Resources {
		if c == subnet.Cidr && n == convertName(subnet.Name) {
			return true, nil
		}
	}

	return false, nil
}
