---
page_title: "jamfpro_disk_encryption_configuration"
description: |-
  
---

# jamfpro_disk_encryption_configuration (Resource)


## Example Usage
```terraform
// jamfpro Institutional Recovery Key config tf example 

resource "jamfpro_disk_encryption_configuration" "disk_encryption_configuration_01" {
  name                     = "jamfpro-tf-example-InstitutionalRecoveryKey-config"
  key_type                 = "Institutional"      # Or "Individual and Institutional"
  file_vault_enabled_users = "Management Account" # Or "Current or Next User"

  institutional_recovery_key {
    certificate_type = "PKCS12" # For .p12 certificate types
    password         = "secretThing"
    data             = filebase64("/Users/dafyddwatkins/localtesting/support_files/filevaultcertificate/FileVaultMaster-sdk.p12")
  }

}

// jamfpro Individual Recovery Key config tf example 

resource "jamfpro_disk_encryption_configuration" "disk_encryption_configuration_02" {
  name                     = "jamfpro-tf-example-IndividualRecoveryKey-config"
  key_type                 = "Individual"
  file_vault_enabled_users = "Management Account" # Or "Current or Next User"

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `file_vault_enabled_users` (String) Defines which user to enable for FileVault 2. Value can be either 'Management Account' or 'Current or Next User'
- `key_type` (String) The type of the key used in the disk encryption which can be either 'Institutional' or 'Individual and Institutional'.
- `name` (String) The name of the disk encryption configuration.

### Optional

- `institutional_recovery_key` (Block List) Details of the institutional recovery key. (see [below for nested schema](#nestedblock--institutional_recovery_key))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The unique identifier of the disk encryption configuration.

<a id="nestedblock--institutional_recovery_key"></a>
### Nested Schema for `institutional_recovery_key`

Optional:

- `certificate_type` (String) The type of certificate used for the institutional recovery key. e.g 'PKCS12' for .p12 certificate types.
- `data` (String) The certificate payload.
- `key` (String)
- `password` (String, Sensitive) The password for the institutional recovery key certificate.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)