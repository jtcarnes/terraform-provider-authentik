package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPPropertyMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLDAPPropertyMappingRead,
		Description: "Get LDAP Property mappings",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"managed_list"},
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"managed_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Retrive multiple property mappings",
			},

			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of ids when `managed_list` is set.",
			},

			"object_field": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLDAPPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PropertymappingsApi.PropertymappingsLdapList(ctx)

	if ml, ok := d.GetOk("managed_list"); ok {
		req = req.Managed(sliceToStringPointer(ml.([]interface{})))
	} else if m, ok := d.GetOk("managed"); ok {
		mm := m.(string)
		req = req.Managed([]*string{&mm})
	}

	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}
	if m, ok := d.GetOk("object_field"); ok {
		req = req.ObjectField(m.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching mappings found")
	}
	if _, ok := d.GetOk("managed_list"); ok {
		d.SetId("-1")
		ids := make([]string, len(res.Results))
		for i, r := range res.Results {
			ids[i] = r.Pk
		}
		d.Set("ids", ids)
	} else {
		f := res.Results[0]
		d.SetId(f.Pk)
		d.Set("name", f.Name)
		d.Set("name", f.Name)
		d.Set("expression", f.Expression)
		d.Set("object_field", f.ObjectField)
	}
	return diags
}
