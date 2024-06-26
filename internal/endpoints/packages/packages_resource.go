// packages_resource.go
package packages

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProPackages defines the schema and CRUD operations for managing Jamf Pro Packages in Terraform.
func ResourceJamfProPackages() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProPackagesCreate,
		ReadContext:   ResourceJamfProPackagesRead,
		UpdateContext: ResourceJamfProPackagesUpdate,
		DeleteContext: ResourceJamfProPackagesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		CustomizeDiff: customValidateFilePath,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the package.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro package.",
			},
			"package_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI of the package in the Jamf Cloud Distribution Service (JCDS).",
			},
			"md5_file_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "md5 hash of the package file for integrity comparison.",
			},
			"package_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file path of the Jamf Pro package.",
			},
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The category of the Jamf Pro package.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok {
						errs = append(errs, fmt.Errorf("%q must be a string, got: %T", key, val))
						return warns, errs
					}

					if v == "" {
						errs = append(errs, fmt.Errorf("%q must not be empty. Either set 'Unknown' to apply no package category or a supply a valid category name string", key))
					} else if v != "Unknown" && len(v) == 0 {
						errs = append(errs, fmt.Errorf("%q must be 'Unknown' or a non-empty string", key))
					}
					return warns, errs
				},
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The filename of the Jamf Pro package.",
			},
			"info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Information about the Jamf Pro package.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes associated with the Jamf Pro package.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The priority of the Jamf Pro package.",
			},
			"reboot_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether a reboot is required after installing the Jamf Pro package.",
			},
			"fill_user_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fill the user template.",
			},
			"fill_existing_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fill existing users.",
			},
			"boot_volume_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether a boot volume is required.",
				Default:     false,
			},
			"allow_uninstalled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to allow the package to be uninstalled.",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OS requirements for the Jamf Pro package.",
			},
			/* Fields are in the data model but don't appear to serve a purpose in jamf 11.3 onwards
			"required_processor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The required processor for the Jamf Pro package.",
				Default:     "None",
			},
			"switch_with_package": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The package to switch with.",
				Default:     "Do Not Install",
			},
			*/
			"install_if_reported_available": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to install the package if it's reported as available.",
			},
			/* Fields are in the data model but don't appear to serve a purpose in jamf 11.3 onwards
			"reinstall_option": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The reinstall option for the Jamf Pro package.",
				Default:     "Do Not Reinstall",
			},
			"triggering_files": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The triggering files for the Jamf Pro package.",
			},
			*/
			"send_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to send a notification for the Jamf Pro package.",
			},
		},
	}
}

// ResourceJamfProPackagesCreate is responsible for creating a new Jamf Pro Package in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProPackagesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize diagnostics
	var diags diag.Diagnostics

	// Extract the file path for the package
	filePath := d.Get("package_file_path").(string)

	// Step 1: Call CreateJCDS2PackageV2 to upload the file to JCDS 2.0
	fileUploadResponse, err := conn.CreateJCDS2PackageV2(filePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to upload file to JCDS 2.0 with file path '%s': %v", filePath, err))
	}
	fmt.Printf("File uploaded successfully, URI: %s\n", fileUploadResponse.URI)

	packageURI := fileUploadResponse.URI

	// After file upload generate the file hash
	fullPath := d.Get("package_file_path").(string)
	fileHash, err := generateMD5FileHash(fullPath)
	if err != nil {
		// Handle error, return diagnostic message
		return diag.FromErr(fmt.Errorf("failed to generate file hash for %s: %v", fullPath, err))
	}

	// Construct the resource object
	packageResourcePointer, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	// Dereference the pointer to get the value
	packageResource := *packageResourcePointer

	// Step 2: Call CreatePackage to create the package metadata in Jamf Pro
	creationResponse, err := conn.CreatePackage(packageResource)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Package '%s': %v", packageResource.Name, err))
	}

	// Set the resource ID, package URI and file hash in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	if err := d.Set("package_uri", packageURI); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("md5_file_hash", fileHash); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetPackageByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Package", strconv.Itoa(creationResponse.ID), checkResourceExists, 45*time.Second, false)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPackagesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPackagesRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPackagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := apiclient.Conn.GetPackageByID(resourceIDInt)

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProPackagesUpdate is responsible for updating an existing Jamf Pro Package on the remote system.
func ResourceJamfProPackagesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics

	// Convert d.Id() from string to integer
	packageID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting package ID '%s' to integer: %v", d.Id(), err))
	}

	// Check if package_file_path has changed
	if d.HasChange("package_file_path") {
		// Step 1: Calculate the new file hash
		filePath := d.Get("package_file_path").(string)
		newFileHash, err := generateMD5FileHash(filePath)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to generate file hash for %s: %v", filePath, err))
		}

		// Step 2: Compare the new file hash with the old one
		oldFileHash, _ := d.GetChange("md5_file_hash")
		if newFileHash != oldFileHash.(string) {
			// The file has changed, upload it
			fileUploadResponse, err := conn.CreateJCDS2PackageV2(filePath)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to upload file to JCDS 2.0 with file path '%s': %v", filePath, err))
			}

			// Update the package_uri and md5_file_hash in Terraform state
			d.Set("package_uri", fileUploadResponse.URI)
			d.Set("md5_file_hash", newFileHash)
			d.Set("filename", filepath.Base(filePath))
		}
	}

	// Update other fields as necessary
	packageResource, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct package for update: %v", err))
	}

	// Update package metadata in Jamf Pro using the integer package ID
	_, err = conn.UpdatePackageByID(packageID, packageResource)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update package with ID %d: %v", packageID, err))
	}

	// Read the updated state
	readDiags := ResourceJamfProPackagesRead(ctx, d, meta)
	if readDiags.HasError() {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPackagesDelete is responsible for deleting a Jamf Pro Package.
func ResourceJamfProPackagesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeletePackageByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeletePackageByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Package '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
