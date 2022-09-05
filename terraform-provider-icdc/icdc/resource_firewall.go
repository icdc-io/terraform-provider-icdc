package icdc

import (
	//"bytes"
	//"encoding/json"
	//"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		Read:   resourceFirewallRead,
		Create: resourceFirewallCreate,
		Update: resourceFirewallUpdate,
		Delete: resourceFirewallDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ems_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ems_ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"source_ip_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_securiry_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceFirewallRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFirewallCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFirewallUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFirewallDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
