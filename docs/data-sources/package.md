---
page_title: "jamfpro_package"
description: |-
  
---

# jamfpro_package (Data Source)


## Example Usage
```terraform
data "jamfpro_package" "jamfpro_package_001_data" {
  id = jamfpro_package.jamfpro_package_001.id
}

output "jamfpro_package_001_data_id" {
  value = data.jamfpro_package.jamfpro_package_001_data.id
}

output "jamfpro_package_001_data_name" {
  value = data.jamfpro_package.jamfpro_package_001_data.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier of the Jamf Pro site.

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `name` (String) The unique name of the Jamf Pro site.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String)