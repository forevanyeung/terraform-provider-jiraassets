resource "jiraassets_object" "example_object" {
  type_id = "100"
  attributes = [
    {
      attr_type_id = "101"
      attr_value   = "My Object"
    },
    {
      attr_type_id = "102"
      attr_value   = "Description of my object"
    }
  ]
}
