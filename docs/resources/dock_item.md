---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jamfpro_dock_item Resource - terraform-provider-jamfpro"
subcategory: ""
description: |-
  
---

# jamfpro_dock_item (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the dock item.
- `path` (String) The path of the dock item.
- `type` (String) The type of the dock item (App/File/Folder).

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `contents` (String) Contents of the dock item.
- `id` (String) The unique identifier of the dock item.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)