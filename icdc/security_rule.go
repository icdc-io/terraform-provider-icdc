package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SecurityRule struct {
	Action          string `json:"action,omitempty"`
	Id              string `json:"id,omitempty"`
	Direction       string `json:"direction,omitempty"`
	PortRangeMin    string `json:"port_range_min,omitempty"`
	PortRangeMax    string `json:"port_range_max,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	NetworkProtocol string `json:"network_protocol,omitempty"`
	RemoteGroupId   string `json:"remote_group_id,omitempty"`
	SecurityGroupId string `json:"security_group_id,omitempty"`
	SourceIpRange   string `json:"source_ip_range,omitempty"`
	EmsRef          string `json:"ems_ref,omitempty"`
}

type MiqTaskDelete struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r *SecurityRule) deleteFromGroup() (bool, error) {

	rb := SecurityRule{
		Action: "remove_firewall_rule",
		Id:     r.EmsRef,
	}

	fmt.Printf("[---DEBUG---] rule %+v", r)

	requestBody, err := json.Marshal(rb)
	requestUrl := fmt.Sprintf("api/compute/v1/security_groups/%s", r.SecurityGroupId)
	responseBody, err := requestApi("POST", requestUrl, bytes.NewBuffer(requestBody))

	var miqTask MiqTaskDelete
	err = responseBody.Decode(&miqTask)

	fmt.Printf("[---DEBUG---] miqTask result %+v", miqTask)

	if err != nil {
		return false, err
	}

	if !miqTask.Success {
		err = fmt.Errorf(miqTask.Message)
		return false, err
	}

	return true, nil
}

func rulesListSnapshot(groupId string) ([]SecurityRule, error) {
	securityGroup, err := fetchSecurityGroup(groupId)

	if err != nil {
		return []SecurityRule{}, err
	}

	return securityGroup.SecurityGroupRules, nil
}
