// advancedmobiledevicesearches_resource.go
package advancedmobiledevicesearches

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"time"

	
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAdvancedMobileDeviceSearches defines the schema for managing mobile device Searches in Terraform.
func ResourceJamfProAdvancedMobileDeviceSearches() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAdvancedMobileDeviceSearchCreate,
		ReadContext:   ResourceJamfProAdvancedMobileDeviceSearchRead,
		UpdateContext: ResourceJamfProAdvancedMobileDeviceSearchUpdate,
		DeleteContext: ResourceJamfProAdvancedMobileDeviceSearchDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the advanced mobile device search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the advanced mobile device search",
			},
			"view_as": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "View type of the mobile device search",
			},
			"sort1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First sorting criteria for the mobile device search",
			},
			"sort2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second sorting criteria for the mobile device search",
			},
			"sort3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Third sorting criteria for the mobile device search",
			},
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"and_or": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"search_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"opening_paren": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"closing_paren": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"display_fields": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "display field in the mobile device search",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "display field item in the mobile device search",
						},
					},
				},
			},
			"site": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID of the site",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the site",
						},
					},
				},
			},
		},
	}
}

// ResourceJamfProAdvancedMobileDeviceSearchCreate is responsible for creating a new Jamf Pro mobile device Search in the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the mobile device search
		search, err := constructJamfProAdvancedMobileDeviceSearch(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the mobile device search for terraform create: %w", err))
		}

		// Log the details of the search that is about to be created
		log.Printf("[INFO] Attempting to create AdvancedMobileDeviceSearch with name: %s", search.Name)

		// Directly call the API to create the resource
		response, err := conn.CreateAdvancedMobileDeviceSearch(search)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating AdvancedMobileDeviceSearch with name: %s. Error: %s", search.Name, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created AdvancedMobileDeviceSearch with ID: %d and name: %s", response.ID, search.Name)

		// Set the ID of the created resource in the Terraform state
		d.SetId(strconv.Itoa(response.ID))

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAdvancedMobileDeviceSearchRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProAdvancedMobileDeviceSearchRead is responsible for reading the current state of a Jamf Pro mobile device Search from the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var search *jamfpro.ResourceAdvancedMobileDeviceSearch

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		searchID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse search ID: %v", convertErr))
		}

		// Try fetching the mobile device search using the ID
		var apiErr error
		search, apiErr = conn.GetAdvancedMobileDeviceSearchByID(searchID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			searchName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			search, apiErr = conn.GetAdvancedMobileDeviceSearchByName(searchName)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Set attributes in the Terraform state
	if err := d.Set("name", search.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("view_as", search.ViewAs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort1", search.Sort1); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort2", search.Sort2); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort3", search.Sort3); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle "criteria" field
	criteriaList := make([]interface{}, len(search.Criteria.Criterion))
	for i, crit := range search.Criteria.Criterion {
		criteriaMap := map[string]interface{}{
			"name":          crit.Name,
			"priority":      crit.Priority,
			"and_or":        crit.AndOr,
			"search_type":   crit.SearchType,
			"value":         crit.Value,
			"opening_paren": crit.OpeningParen,
			"closing_paren": crit.ClosingParen,
		}
		criteriaList[i] = criteriaMap
	}
	if err := d.Set("criteria", criteriaList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle "display_fields" field
	if len(search.DisplayFields) == 0 || len(search.DisplayFields[0].DisplayField) == 0 {
		if err := d.Set("display_fields", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		displayFieldsList := make([]map[string]interface{}, len(search.DisplayFields[0].DisplayField))
		for i, displayField := range search.DisplayFields[0].DisplayField {
			displayFieldMap := map[string]interface{}{
				"name": displayField.Name,
			}
			displayFieldsList[i] = displayFieldMap
		}
		if err := d.Set("display_fields", displayFieldsList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Handle "site" field
	site := map[string]interface{}{
		"id":   search.Site.ID,
		"name": search.Site.Name,
	}
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProAdvancedMobileDeviceSearchUpdate is responsible for updating an existing Jamf Pro mobile device Search on the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated mobile device search
		search, err := constructJamfProAdvancedMobileDeviceSearch(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the mobile device search for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		searchID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse search ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateAdvancedMobileDeviceSearchByID(searchID, search)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			searchName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
			_, apiErr = conn.UpdateAdvancedMobileDeviceSearchByName(searchName, search)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAdvancedMobileDeviceSearchRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProAdvancedMobileDeviceSearchDelete is responsible for deleting a Jamf Pro AdvancedMobileDeviceSearch.
func ResourceJamfProAdvancedMobileDeviceSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		groupID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse group ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteAdvancedMobileDeviceSearchByID(groupID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
			apiErr = conn.DeleteAdvancedMobileDeviceSearchByName(groupName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
