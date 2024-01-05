package icdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   resourceSecurityGroupRead,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,
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
				Type:                  schema.TypeString,
				Required:              true,
				DiffSuppressOnRefresh: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return convertName(old) == convertName(new)
				},
			},
		},
	}
}

func resourceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	defer ctx.Done()

	emsId, err := fetchEmsId()

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	snapshot, err := groupsListSnapshot()

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	securityGroupCreateRequest := SecurityGroupCreateRequest{
		Action: "create",
		Name:   convertName(d.Get("name").(string)),
	}

	requestBody, err := json.Marshal(securityGroupCreateRequest)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	requestUrl := fmt.Sprintf("api/compute/v1/providers/%s/security_groups", emsId)
	responseBody, err := requestApi("POST", requestUrl, bytes.NewBuffer(requestBody))

	fmt.Printf("[---DEBUG--] responseBody %+v", responseBody)

	var miqTaskResults MiqTaskResults
	err = responseBody.Decode(&miqTaskResults)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = retry.RetryContext(
		ctx,
		d.Timeout(schema.TimeoutCreate),
		func() *retry.RetryError {
			requestUrl := fmt.Sprintf("api/compute/v1/tasks/%s", miqTaskResults.Results[0].TaskId)
			responseBody, err := requestApi("GET", requestUrl, nil)

			if err != nil {
				return retry.NonRetryableError(err)
			}

			var miqTask MiqTask
			err = responseBody.Decode(&miqTask)

			if err != nil {
				return retry.NonRetryableError(err)
			}

			if miqTask.State == "Finished" && miqTask.Status == "Ok" {
				return nil
			}

			return retry.RetryableError(fmt.Errorf("waiting for creating security group task finished: %+v", miqTask))
		},
	)

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	groups, err := securityGroupList()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	var createdGroup SecurityGroup

	for _, g := range groups {
		_, ok := snapshot[g.Id]

		if !ok {
			createdGroup = g
			break
		}
	}

	fmt.Printf("[---DEBUG---] createdGroup %+v", createdGroup)

	err = retry.RetryContext(
		ctx,
		d.Timeout(schema.TimeoutCreate),
		func() *retry.RetryError {
			fmt.Println("[---DEBUG---] waiting for security group rules")

			if len(createdGroup.SecurityGroupRules) > 0 {
				return nil
			}

			createdGroup, err = fetchSecurityGroup(createdGroup.Id)

			if err != nil {
				if err != nil {
					return retry.NonRetryableError(err)
				}
			}

			return retry.RetryableError(fmt.Errorf("waiting for creating security group (%s) rules", createdGroup.Id))
		},
	)

	err = d.Set("ems_ref", createdGroup.EmsRef)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("name", createdGroup.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(createdGroup.Id)

	return nil
}

func resourceSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
