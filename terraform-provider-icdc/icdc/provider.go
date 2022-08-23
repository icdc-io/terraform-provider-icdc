package icdc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"icdc_service_request": resourceServiceRequest(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
