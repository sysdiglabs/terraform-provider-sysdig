package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureUserRulesFile() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resouceSysdigUserRulesFileUpdate,
		Read:   resouceSysdigUserRulesFileRead,
		Update: resouceSysdigUserRulesFileUpdate,
		Delete: resouceSysdigUserRulesFileDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resouceSysdigUserRulesFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	userRulesFile, err := updateCurrentUserRulesFile(client, d.Get("content").(string))
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(userRulesFile.Version))
	d.Set("version", userRulesFile.Version)

	return nil
}

func updateCurrentUserRulesFile(client secure.SysdigSecureClient, content string) (secure.UserRulesFile, error) {
	userRulesFile, err := client.GetUserRulesFile()
	if err != nil {
		return secure.UserRulesFile{}, err
	}
	userRulesFile.Content = content

	userRulesFile, err = client.UpdateUserRulesFile(userRulesFile)
	if err != nil {
		return secure.UserRulesFile{}, err
	}

	return userRulesFile, nil
}

func resouceSysdigUserRulesFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	userRulesFile, err := client.GetUserRulesFile()
	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("content", userRulesFile.Content)
	d.Set("version", userRulesFile.Version)

	return nil
}

func resouceSysdigUserRulesFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	_, err := updateCurrentUserRulesFile(client, "")

	return err
}
