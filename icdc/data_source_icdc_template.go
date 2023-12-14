package icdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceICDCTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceICDCTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type ICDCTemplate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func dataSourceICDCTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filter := d.Get("name").(string)
	templates, err := fetchTemplatesByFilter(filter)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(templates) == 0 {
		return diag.Errorf("No templates found for the provided filter")
	}

	if err := d.Set("name", templates[0].Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(templates[0].Id)

	return nil
}

type ICDCTemplateCollection struct {
	Resources []ICDCTemplate `json:"resources"`
}

func fetchTemplatesByFilter(filter string) ([]ICDCTemplate, error) {
	requestUrl := "api/compute/v1/service_templates?expand=resources&filter[]=name='%25" + filter + "%25'"
	responseBody, err := requestApi("GET", requestUrl, nil)

	if err != nil {
		return nil, err
	}

	var icdcTemplateCollection *ICDCTemplateCollection
	responseBody.Decode(&icdcTemplateCollection)

	return icdcTemplateCollection.Resources, nil
}
