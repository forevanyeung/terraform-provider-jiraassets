---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jiraassets Provider"
subcategory: ""
description: |-
  A Terraform provider for Jira Assets.
---

# jiraassets Provider

A Terraform provider for Jira Assets.

## Example Usage

```terraform
provider "jiraassets" {
  workspace_id = ""
  user         = ""
  password     = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `password` (String, Sensitive) Personal access token for the admin or service account.
- `user` (String) Username of an admin or service account with access to the Jira API.
- `workspace_id` (String) Workspace Id of the Assets instance.
