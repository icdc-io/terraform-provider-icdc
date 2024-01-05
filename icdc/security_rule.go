package icdc

type SecurityRule struct {
	Action          string `json:"action,omitempty"`
	Id              string `json:"id,omitempty"`
	Direction       string `json:"direction"`
	PortRangeMin    string `json:"port_range_min"`
	PortRangeMax    string `json:"port_range_max"`
	Protocol        string `json:"protocol"`
	NetworkProtocol string `json:"network_protocol"`
	RemoteGroupId   string `json:"remote_group_id"`
	SecurityGroupId string `json:"security_group_id"`
	SourceIpRange   string `json:"source_ip_range"`
}

func (r *SecurityRule) deleteFromGroup() (bool, error) {
	return true, nil
}

func rulesListSnapshot(groupId string) ([]SecurityRule, error) {
	securityGroup, err := fetchSecurityGroup(groupId)

	if err != nil {
		return []SecurityRule{}, err
	}

	return securityGroup.SecurityGroupRules, nil
}
