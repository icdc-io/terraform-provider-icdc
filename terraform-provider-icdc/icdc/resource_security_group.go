package icdc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSecurityGroupRead,
		Create: resourceSecurityGroupCreate,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,
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
			"egress": {
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
			"ingress": {
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

func resourceSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
