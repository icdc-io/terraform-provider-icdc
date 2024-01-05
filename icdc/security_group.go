package icdc

import (
	"fmt"
)

type SecurityGroup struct {
	Id                 string         `json:"id,omitempty"`
	Name               string         `json:"name"`
	EmsRef             string         `json:"ems_ref"`
	SecurityGroupRules []SecurityRule `json:"firewall_rules,omitempty"`
}

type SecurityGroupCollection struct {
	Resources []SecurityGroup `json:"resources"`
}

type SecurityGroupCreateRequest struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

type MiqTaskResults struct {
	Results []struct {
		TaskId   string `json:"task_id"`
		Success  bool   `json:"success"`
		TaskHref string `json:"task_href"`
	} `json:"results"`
}

type MiqTask struct {
	Id      string `json:"id"`
	State   string `json:"state"`
	Status  string `json:"status"`
	Success string `json:"success"`
	Message string `json:"message"`
}

type EmsProviderCollection struct {
	Resources []struct {
		Id string `json:"id"`
	}
}

func fetchEmsId() (string, error) {

	requestUrl := "api/compute/v1/providers?expand=resources&filter[]=type=ManageIQ::Providers::Redhat::NetworkManager"

	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return "", err
	}

	var parsedBody EmsProviderCollection

	err = responseBody.Decode(&parsedBody)

	if err != nil {
		return "", err
	}

	return parsedBody.Resources[0].Id, nil
}

func groupsListSnapshot() (map[string]string, error) {
	securityGroups, err := securityGroupList()

	if err != nil {
		return nil, err
	}

	snapshot := make(map[string]string)

	for _, s := range securityGroups {
		snapshot[s.Id] = s.Name
	}

	return snapshot, nil
}

func securityGroupList() ([]SecurityGroup, error) {
	requestUrl := "api/compute/v1/security_groups?expand=resources&attributes=firewall_rules"
	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return nil, fmt.Errorf("can't fetch security group list: %s", err)
	}

	var securityGroupCollection SecurityGroupCollection

	err = responseBody.Decode(&securityGroupCollection)

	if err != nil {
		return nil, fmt.Errorf("can't decode security group list: %s", err)
	}

	return securityGroupCollection.Resources, nil
}

func fetchSecurityGroup(id string) (SecurityGroup, error) {
	requestUrl := fmt.Sprintf("api/compute/v1/security_groups/%s?expand=resources&attributes=firewall_rules", id)

	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return SecurityGroup{}, fmt.Errorf("can't fetch security group: %s", err)
	}

	var securityGroup SecurityGroup

	err = responseBody.Decode(&securityGroup)

	if err != nil {
		return SecurityGroup{}, fmt.Errorf("can't decode security group: %s", err)
	}

	return securityGroup, nil
}
