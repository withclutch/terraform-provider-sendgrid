/*
Provide a resource to manage a domain authentication validation.
Example Usage
```hcl

	resource "sendgrid_domain_authentication_validation" "foo" {
		domain_authentication_id = sendgrid_domain_authentication.foo.id
	}

```
*/
package sendgrid

import (
	"context"
	"fmt"

	sendgrid "github.com/arslanbekov/terraform-provider-sendgrid/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://docs.sendgrid.com/api-reference/domain-authentication/validate-a-domain-authentication
func resourceSendgridDomainAuthenticationValidation() *schema.Resource { //nolint:funlen
	return &schema.Resource{
		CreateContext: resourceSendgridDomainAuthenticationValidationCreate,
		ReadContext:   resourceSendgridDomainAuthenticationValidationRead,
		DeleteContext: resourceSendgridDomainAuthenticationValidationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_authentication_id": {
				Type:        schema.TypeString,
				Description: "Id of the domain authentication to validate.",
				Required:    true,
				ForceNew:    true,
			},

			"valid": {
				Type:        schema.TypeBool,
				Description: "Indicates if this is a valid authenticated domain or not.",
				Computed:    true,
			},
		},
	}
}

func resourceSendgridDomainAuthenticationValidationCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return validateDomain(ctx, d, m)
}

func validateDomain(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	if err := c.ValidateDomainAuthentication(ctx, d.Get("domain_authentication_id").(string)); err.Err != nil || err.StatusCode != 200 {
		if err.Err != nil {
			return diag.FromErr(err.Err)
		}
		return diag.Errorf("unable to validate domain DNS configuration")
	}

	d.SetId(d.Get("domain_authentication_id").(string))

	return resourceSendgridDomainAuthenticationValidationRead(ctx, d, m)
}

func resourceSendgridDomainAuthenticationValidationRead( //nolint:funlen,cyclop
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	auth, err := c.ReadDomainAuthentication(ctx, d.Get("domain_authentication_id").(string))
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	if err := d.Set("valid", auth.Valid); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(auth.ID))
	return nil
}

func resourceSendgridDomainAuthenticationValidationDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
