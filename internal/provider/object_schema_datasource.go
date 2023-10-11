package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewObjectSchemaDataSource() datasource.DataSource {
    return &objectSchemaDataSource{}
}

type objectSchemaDataSource struct{}

func (d *objectSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_object_schema"
}

func (d *objectSchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

func (d *objectSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
