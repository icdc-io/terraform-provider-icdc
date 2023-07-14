package icdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		Read:   resourceVpcRead,
		Create: resourceVpcCreate,
		Update: resourceVpcUpdate,
		Delete: resourceVpcDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"router": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVpcRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s", d.Get("id"))
	r, err := requestApi("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error getting api services: %w", err)
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var vpcRequestResponse *VpcRequestResponse

	if err = json.Unmarshal(resBody, &vpcRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	log.Println(PrettyStruct(vpcRequestResponse))
	log.Println(vpcRequestResponse.Vpc.Id)
	d.SetId(vpcRequestResponse.Vpc.Id)
	return nil
}

func resourceVpcCreate(d *schema.ResourceData, m interface{}) error {
	cloudVpcRaw := &VpcCreateBody{
		Vpc: VpcStructBody{
			Name: d.Get("name").(string),
			Router: RouterCreateBody{
				Name: d.Get("name").(string),
			},
		},
	}

	requestBody, err := json.Marshal(cloudVpcRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudVpcRaw))

	url := ""
	r, err := requestApi("POST", url, body)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var vpcRequestResponse *VpcRequestResponse

	if err = json.Unmarshal(resBody, &vpcRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(vpcRequestResponse))
	log.Println(PrettyStruct(vpcRequestResponse))

	vpc_id := vpcRequestResponse.Vpc.Id
	log.Println(PrettyStruct(vpc_id))
	d.SetId(vpc_id)
	return nil
}

func resourceVpcUpdate(d *schema.ResourceData, m interface{}) error {

	cloudVpcRaw := &VpcCreateBody{
		Vpc: VpcStructBody{
			Name: d.Get("name").(string),
			Router: RouterCreateBody{
				Name: d.Get("name").(string),
			},
		},
	}

	requestBody, err := json.Marshal(cloudVpcRaw)
	if err != nil {
		return fmt.Errorf("error marshaling service request: %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	log.Println(PrettyStruct(cloudVpcRaw))

	url := fmt.Sprintf("%s", d.Get("id").(string))
	r, err := requestApi("PUT", url, body)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var vpcRequestResponse *VpcRequestResponse

	if err = json.Unmarshal(resBody, &vpcRequestResponse); err != nil {
		return fmt.Errorf("error decoding service response: %w", err)
	}

	fmt.Println(PrettyStruct(vpcRequestResponse))
	vpc_id := vpcRequestResponse.Vpc.Id
	d.SetId(vpc_id)
	return nil
}

func resourceVpcDelete(d *schema.ResourceData, m interface{}) error {

	url := fmt.Sprintf("%s", d.Get("id").(string))
	_, err := requestApi("DELETE", url, nil)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil

}
