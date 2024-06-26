---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jamfpro_macos_configuration_profile Resource - terraform-provider-jamfpro"
subcategory: ""
description: |-
  
---

# jamfpro_macos_configuration_profile (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Jamf UI name for configuration profile.
- `payload` (String) A MacOS configuration profile xml file as a file
- `scope` (Block List, Min: 1, Max: 1) The scope of the configuration profile. (see [below for nested schema](#nestedblock--scope))

### Optional

- `category` (Block List, Max: 1) The category to which the configuration profile is scoped. (see [below for nested schema](#nestedblock--category))
- `description` (String) Description of the configuration profile.
- `distribution_method` (String) The distribution method for the configuration profile. ['Make Available in Self Service','Install Automatically']
- `level` (String) The level of the configuration profile. Available options are: 'Computer', 'User' or 'System'.
- `redeploy_on_update` (String) Whether the configuration profile is redeployed on update.
- `self_service` (Block List, Max: 1) Self Service Configuration (see [below for nested schema](#nestedblock--self_service))
- `site` (Block List, Max: 1) The site to which the configuration profile is scoped. (see [below for nested schema](#nestedblock--site))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_removeable` (Boolean) Whether the configuration profile is user removeable or not.

### Read-Only

- `id` (String) The unique identifier of the macOS configuration profile.
- `uuid` (String) The UUID of the configuration profile.

<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Required:

- `all_computers` (Boolean) Whether the configuration profile is scoped to all computers.

Optional:

- `all_jss_users` (Boolean) Whether the configuration profile is scoped to all JSS users.
- `building_ids` (List of Number) The buildings to which the configuration profile is scoped by Jamf ID
- `computer_group_ids` (List of Number) The computer groups to which the configuration profile is scoped by Jamf ID
- `computer_ids` (List of Number) The computers to which the configuration profile is scoped by Jamf ID
- `department_ids` (List of Number) The departments to which the configuration profile is scoped by Jamf ID
- `exclusions` (Block List, Max: 1) The exclusions from the scope. (see [below for nested schema](#nestedblock--scope--exclusions))
- `jss_user_group_ids` (List of Number) The jss user groups to which the configuration profile is scoped by Jamf ID
- `jss_user_ids` (List of Number) The jss users to which the configuration profile is scoped by Jamf ID
- `limitations` (Block List, Max: 1) The limitations within the scope. (see [below for nested schema](#nestedblock--scope--limitations))

<a id="nestedblock--scope--exclusions"></a>
### Nested Schema for `scope.exclusions`

Optional:

- `building_ids` (List of Number) Buildings excluded from scope by Jamf ID.
- `computer_group_ids` (List of Number) Computer Groups excluded from scope by Jamf ID.
- `computer_ids` (List of Number) Computers excluded from scope by Jamf ID.
- `department_ids` (List of Number) Departments excluded from scope by Jamf ID.
- `ibeacon_ids` (List of Number) Ibeacons excluded from scope by Jamf ID.
- `jss_user_group_ids` (List of Number) JSS User Groups excluded from scope by Jamf ID.
- `jss_user_ids` (List of Number) JSS Users excluded from scope by Jamf ID.
- `network_segment_ids` (List of Number) Network segments excluded from scope by Jamf ID.


<a id="nestedblock--scope--limitations"></a>
### Nested Schema for `scope.limitations`

Optional:

- `ibeacon_ids` (List of Number) Ibeacons the scope is limited to by Jamf ID.
- `network_segment_ids` (List of Number) Network segments the scope is limited to by Jamf ID.
- `user_group_ids` (List of Number) Users groups the scope is limited to by Jamf ID.
- `user_names` (List of String) Users the macOS config profile scope is limited to by Jamf ID.



<a id="nestedblock--category"></a>
### Nested Schema for `category`

Required:

- `id` (Number) The unique identifier of the category to which the configuration profile is scoped.

Optional:

- `name` (String) The name of the category to which the configuration profile is scoped.


<a id="nestedblock--self_service"></a>
### Nested Schema for `self_service`

Optional:

- `feature_on_main_page` (Boolean) Shows Configuration Profile on Self Service main page
- `force_users_to_view_description` (Boolean) Forces users to view the description
- `install_button_text` (String) Text shown on Self Service install button
- `notification` (Boolean) Enables Notification for this profile in self service
- `notification_message` (String) Message body
- `notification_subject` (String) Message Subject
- `self_service_categories` (Block List) Self Service category options (see [below for nested schema](#nestedblock--self_service--self_service_categories))
- `self_service_description` (String) Description shown in Self Service

<a id="nestedblock--self_service--self_service_categories"></a>
### Nested Schema for `self_service.self_service_categories`

Required:

- `display_in` (Boolean) Display this profile in this category?
- `feature_in` (Boolean) Feature this profile in this category?

Optional:

- `id` (Number) ID of category
- `name` (String) Name of category



<a id="nestedblock--site"></a>
### Nested Schema for `site`

Optional:

- `id` (Number) Jamf Pro Site ID. Value defaults to -1 aka not used.
- `name` (String) Jamf Pro Site Name. Value defaults to 'None' aka not used


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)
