package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
)

var (
	_ provider.Provider = &awsmtProvider{}
)

func New() provider.Provider {
	return &awsmtProvider{}
}

type awsmtProvider struct{}

type awsmtProviderModel struct {
	Profile types.String `tfsdk:"profile"`
	Region  types.String `tfsdk:"region"`
}

func (p *awsmtProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awsmt"
}

func (p *awsmtProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "The profile generated by the SSO login. You can find the profile(s) name in '~/.aws/config'. SSO login will not be used if the profile name is not specified and no environmental variable called 'aws_profile' is found.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "AWS region. defaults to 'eu-central-1'.",
			},
		},
	}
}

func (p *awsmtProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring AWS MediaTailor client")

	var config awsmtProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var region = "eu-central-1"
	var profile = ""

	if !config.Region.IsUnknown() || !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	var sess *session.Session
	var err error

	if config.Profile.IsUnknown() || config.Profile.IsNull() || config.Profile.ValueString() == "" {
		if os.Getenv("AWS_PROFILE") != "" {
			profile = os.Getenv("AWS_PROFILE")
		}
	} else {
		profile = config.Profile.ValueString()
	}

	tflog.Debug(ctx, "Creating AWS client session")

	if profile != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: aws.String(region),
			},
			Profile: profile,
		})
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to Initialize Provider in Region", "unable to initialize provider in the specified region: "+err.Error())
		return
	}

	c := mediatailor.New(sess)

	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "AWS MediaTailor client configured", map[string]any{"success": true})
}

func (p *awsmtProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		DataSourceChannel,
		DataSourceSourceLocation,
		DataSourcePlaybackConfiguration,
		DataSourceLiveSource,
		DataSourceVodSource,
	}

}

func (p *awsmtProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		ResourceChannel,
		ResourceSourceLocation,
		ResourcePlaybackConfiguration,
		ResourceLiveSource,
		ResourceVodSource,
	}
}
