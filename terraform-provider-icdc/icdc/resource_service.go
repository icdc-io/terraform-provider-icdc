package icdc
/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const baseURL = "https://api.ycz.icdc.io/api/compute/v1"
const miqGroup = "icdc.member"



func resourceServiceRequestRead(d *schema.ResourceData, m interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/services?expand=resources", baseURL))

	req.Header.set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")))
	req.Header.Set("X_MIQ_GROUP", miqGroup)

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	var response *ServiceRequestResponse

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return err
	}


	return nil
}

func resourceServiceRequestCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceRequestUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceRequestDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceRequest() *schema.Resource {
	return &schema.Resource{
		Read: resourceServiceRead,
		Create: resourceServiceCreate,
		Update: resourceServiceUpdate,
		Delete: resourceServiceDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"href": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"guid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
*/