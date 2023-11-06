---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jiraassets_object_schema Data Source - terraform-provider-jiraassets"
subcategory: ""
description: |-
  
---

# jiraassets_object_schema (Data Source)



## Example Usage

```terraform
data "jiraassets_object_schema" "example" {
  id = "100"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of the object schema.

### Read-Only

- `can_manage` (Boolean)
- `created` (String)
- `description` (String)
- `global_id` (String)
- `id_as_int` (Number)
- `name` (String)
- `object_count` (Number)
- `object_schema_key` (String)
- `object_type_count` (Number)
- `status` (String)
- `updated` (String)
- `workspace_id` (String)