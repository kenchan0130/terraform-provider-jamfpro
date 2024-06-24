// categories_resource.go
package categories

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProCategories defines the schema and CRUD operations for managing Jamf Pro Categories in Terraform.
func ResourceJamfProCategories() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProCategoriesCreate,
		ReadContext:   resourceJamfProCategoriesReadWithCleanup,
		UpdateContext: resourceJamfProCategoriesUpdate,
		DeleteContext: resourceJamfProCategoriesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(15 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the category.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro category.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     9,
				Description: "The priority of the Jamf Pro category.",
			},
		},
	}
}
