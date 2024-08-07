---
page_title: "jamfpro_restricted_software"
description: |-
  
---

# jamfpro_restricted_software (Data Source)


## Example Usage
```terraform
data "jamfpro_restricted_software" "restricted_software_001_data" {
  id = jamfpro_restricted_software.restricted_software_001.id
}

output "jamfpro_restricted_software_001_id" {
  value = data.jamfpro_restricted_software.restricted_software_001_data.id
}

output "jamfpro_restricted_software_001_name" {
  value = data.jamfpro_restricted_software.restricted_software_001_data.name
}

data "jamfpro_restricted_software" "restricted_software_002_data" {
  id = jamfpro_restricted_software.restricted_software_002.id
}

output "jamfpro_restricted_software_002_id" {
  value = data.jamfpro_restricted_software.restricted_software_002_data.id
}

output "jamfpro_restricted_software_002_name" {
  value = data.jamfpro_restricted_software.restricted_software_002_data.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier of the Jamf Pro restricted software.

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `name` (String) The unique name of the Jamf Pro restricted software.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String)