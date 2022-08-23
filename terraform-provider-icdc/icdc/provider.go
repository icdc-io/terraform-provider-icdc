package icdc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"icdc_service_request": resourceServiceRequest(),
			"icdc_service": resourceService(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
