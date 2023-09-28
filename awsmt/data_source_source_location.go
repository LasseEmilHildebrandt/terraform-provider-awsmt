package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ datasource.DataSource              = &dataSourceSourceLocation{}
	_ datasource.DataSourceWithConfigure = &dataSourceSourceLocation{}
)

func DataSourceSourceLocation() datasource.DataSource {
	return &dataSourceSourceLocation{}
}

type dataSourceSourceLocation struct {
	client *mediatailor.MediaTailor
}

func (d *dataSourceSourceLocation) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
}

func (d *dataSourceSourceLocation) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": computedString,
			"access_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"access_type": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("S3_SIGV4", "SECRETS_MANAGER_ACCESS_TOKEN"),
						},
					},
					"secrets_manager_access_token_configuration": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"header_name":       computedString,
							"secret_arn":        computedString,
							"secret_string_key": computedString,
						},
					},
				},
			},
			"arn":           computedString,
			"creation_time": computedString,
			"default_segment_delivery_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"dsdc_base_url": computedString,
				},
			},
			"http_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"hc_base_url": computedString,
				},
			},
			"last_modified_time": computedString,
			"segment_delivery_configurations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sdc_base_url": computedString,
						"sdc_name":     computedString,
					},
				},
			},
			"source_location_name": requiredString,
			"tags":                 computedMap,
		},
	}
}

func (d *dataSourceSourceLocation) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceSourceLocation) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sourceLocationModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := data.SourceLocationName

	sourceLocation, err := d.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: sourceLocationName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", err.Error())
		return
	}

	data = readSourceLocationToPlan(data, mediatailor.CreateSourceLocationOutput(*sourceLocation))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
