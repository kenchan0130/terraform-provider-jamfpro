// mobiledeviceconfigurationprofilesplist_state.go
package mobiledeviceconfigurationprofilesplist

import (
	"sort"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest ResourceMobileDeviceConfigurationProfile
// information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceMobileDeviceConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a map to hold the resource data
	resourceData := map[string]interface{}{
		"name":              resource.General.Name,
		"description":       resource.General.Description,
		"uuid":              resource.General.UUID,
		"deployment_method": resource.General.DeploymentMethod,
		// Skipping the 'distribution_method' attribute as it appears to be deprecated but still in documentation
		"redeploy_on_update":                resource.General.RedeployOnUpdate,
		"redeploy_days_before_cert_expires": resource.General.RedeployDaysBeforeCertExpires,
	}

	// Check if the level is "System" and set it to "Device Level", otherwise use the value from resource
	// This is done to match the Jamf Pro API behavior
	levelValue := resource.General.Level
	if levelValue == "System" {
		levelValue = "Device Level"
	}
	resourceData["level"] = levelValue

	// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
	site := []interface{}{}
	if resource.General.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id": resource.General.Site.ID,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Sanitize and set the payloads using the plist processor function
	profile := plist.NormalizePayloadState(resource.General.Payloads)
	if err := d.Set("payloads", profile); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'category' attribute in the state only if it's not empty (i.e., not default values)
	category := []interface{}{}
	if resource.General.Category.ID != -1 {
		category = append(category, map[string]interface{}{
			"id": resource.General.Category.ID,
		})
	}
	if len(category) > 0 {
		if err := d.Set("category", category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Preparing and setting scope data
	if scopeData, err := setScope(resource); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else if err := d.Set("scope", []interface{}{scopeData}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update the resource data
	for k, v := range resourceData {
		if err := d.Set(k, v); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setScope converts the scope structure into a format suitable for setting in the Terraform state.
func setScope(resource *jamfpro.ResourceMobileDeviceConfigurationProfile) (map[string]interface{}, error) {
	scopeData := map[string]interface{}{
		"all_mobile_devices": resource.Scope.AllMobileDevices,
		"all_jss_users":      resource.Scope.AllJSSUsers,
	}

	// Gather mobile devices, groups, etc.
	scopeData["mobile_device_ids"] = flattenAndSortMobileDeviceIDs(resource.Scope.MobileDevices)
	scopeData["mobile_device_group_ids"] = flattenAndSortScopeEntityIds(resource.Scope.MobileDeviceGroups)
	scopeData["jss_user_ids"] = flattenAndSortScopeEntityIds(resource.Scope.JSSUsers)
	scopeData["jss_user_group_ids"] = flattenAndSortScopeEntityIds(resource.Scope.JSSUserGroups)
	scopeData["building_ids"] = flattenAndSortScopeEntityIds(resource.Scope.Buildings)
	scopeData["department_ids"] = flattenAndSortScopeEntityIds(resource.Scope.Departments)

	// Gather limitations
	limitationsData, err := setLimitations(resource.Scope.Limitations)
	if err != nil {
		return nil, err
	}
	if limitationsData != nil {
		scopeData["limitations"] = limitationsData
	}

	// Gather exclusions
	exclusionsData, err := setExclusions(resource.Scope.Exclusions)
	if err != nil {
		return nil, err
	}
	if exclusionsData != nil {
		scopeData["exclusions"] = exclusionsData
	}

	return scopeData, nil
}

// setLimitations collects and formats limitations data for the Terraform state.
func setLimitations(limitations jamfpro.MobileDeviceConfigurationProfileSubsetLimitation) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(limitations.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIds(limitations.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(limitations.Ibeacons) > 0 {
		ibeaconIDs := flattenAndSortScopeEntityIds(limitations.Ibeacons)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
		}
	}

	if len(limitations.Users) > 0 {
		userNames := flattenAndSortScopeEntityNames(limitations.Users)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(limitations.UserGroups) > 0 {
		userGroupIDs := flattenAndSortScopeEntityIds(limitations.UserGroups)
		if len(userGroupIDs) > 0 {
			result["directory_service_usergroup_ids"] = userGroupIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// setExclusions collects and formats exclusion data for the Terraform state.
func setExclusions(exclusions jamfpro.MobileDeviceConfigurationProfileSubsetExclusion) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(exclusions.MobileDevices) > 0 {
		computerIDs := flattenAndSortMobileDeviceIDs(exclusions.MobileDevices)
		if len(computerIDs) > 0 {
			result["mobile_device_ids"] = computerIDs
		}
	}

	if len(exclusions.MobileDeviceGroups) > 0 {
		computerGroupIDs := flattenAndSortScopeEntityIds(exclusions.MobileDeviceGroups)
		if len(computerGroupIDs) > 0 {
			result["mobile_device_group_ids"] = computerGroupIDs
		}
	}

	if len(exclusions.Buildings) > 0 {
		buildingIDs := flattenAndSortScopeEntityIds(exclusions.Buildings)
		if len(buildingIDs) > 0 {
			result["building_ids"] = buildingIDs
		}
	}

	if len(exclusions.JSSUsers) > 0 {
		jssUserIDs := flattenAndSortScopeEntityIds(exclusions.JSSUsers)
		if len(jssUserIDs) > 0 {
			result["jss_user_ids"] = jssUserIDs
		}
	}

	if len(exclusions.JSSUserGroups) > 0 {
		jssUserGroupIDs := flattenAndSortScopeEntityIds(exclusions.JSSUserGroups)
		if len(jssUserGroupIDs) > 0 {
			result["jss_user_group_ids"] = jssUserGroupIDs
		}
	}

	if len(exclusions.Departments) > 0 {
		departmentIDs := flattenAndSortScopeEntityIds(exclusions.Departments)
		if len(departmentIDs) > 0 {
			result["department_ids"] = departmentIDs
		}
	}

	if len(exclusions.NetworkSegments) > 0 {
		networkSegmentIDs := flattenAndSortNetworkSegmentIds(exclusions.NetworkSegments)
		if len(networkSegmentIDs) > 0 {
			result["network_segment_ids"] = networkSegmentIDs
		}
	}

	if len(exclusions.Users) > 0 {
		userNames := flattenAndSortScopeEntityNames(exclusions.Users)
		if len(userNames) > 0 {
			result["directory_service_or_local_usernames"] = userNames
		}
	}

	if len(exclusions.UserGroups) > 0 {
		userGroupIDs := flattenAndSortScopeEntityIds(exclusions.UserGroups)
		if len(userGroupIDs) > 0 {
			result["directory_service_usergroup_ids"] = userGroupIDs
		}
	}

	if len(exclusions.IBeacons) > 0 {
		ibeaconIDs := flattenAndSortScopeEntityIds(exclusions.IBeacons)
		if len(ibeaconIDs) > 0 {
			result["ibeacon_ids"] = ibeaconIDs
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return []map[string]interface{}{result}, nil
}

// helper functions

// flattenAndSortScopeEntityIds converts a slice of general scope entities (like user groups, buildings) to a format suitable for Terraform state.
func flattenAndSortScopeEntityIds(entities []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) []int {
	var ids []int
	for _, entity := range entities {
		if entity.ID != 0 {
			ids = append(ids, entity.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortScopeEntityNames converts a slice of RestrictedSoftwareSubsetScopeEntity into a sorted slice of strings.
func flattenAndSortScopeEntityNames(entities []jamfpro.MobileDeviceConfigurationProfileSubsetScopeEntity) []string {
	var names []string
	for _, entity := range entities {
		if entity.Name != "" {
			names = append(names, entity.Name)
		}
	}
	sort.Strings(names)
	return names
}

// flattenAndSortMobileDeviceIDs converts a slice of MobileDeviceConfigurationProfileSubsetMobileDevice into a sorted slice of integers.
func flattenAndSortMobileDeviceIDs(devices []jamfpro.MobileDeviceConfigurationProfileSubsetMobileDevice) []int {
	var ids []int
	for _, device := range devices {
		if device.ID != 0 {
			ids = append(ids, device.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

// flattenAndSortNetworkSegmentIds converts a slice of MobileDeviceConfigurationProfileSubsetNetworkSegment into a sorted slice of integers.
func flattenAndSortNetworkSegmentIds(segments []jamfpro.MobileDeviceConfigurationProfileSubsetNetworkSegment) []int {
	var ids []int
	for _, segment := range segments {
		if segment.ID != 0 {
			ids = append(ids, segment.ID)
		}
	}
	sort.Ints(ids)
	return ids
}