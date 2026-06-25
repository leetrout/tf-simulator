package provider

import (
	"context"

	"github.com/leetrout/terraform-sim/cmd/terraform-provider-nimbus/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultEndpoint = "http://localhost:9321"

// Ensure NimbusProvider satisfies the provider.Provider interface.
var _ provider.Provider = (*NimbusProvider)(nil)

// NimbusProvider is the Terraform provider for the tfsim cloud API.
type NimbusProvider struct {
	// version is set by the build (e.g. via ldflags) and surfaced to Terraform.
	version string
}

// New returns a provider.Provider factory for the given version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NimbusProvider{version: version}
	}
}

// nimbusProviderModel maps provider configuration.
type nimbusProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *NimbusProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nimbus"
	resp.Version = p.version
}

func (p *NimbusProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The nimbus provider manages resources in a running tfsim cloud (HTTP API at /api/cloud/*).",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Base URL of the tfsim server. Defaults to " + defaultEndpoint + ".",
				Optional:            true,
			},
		},
	}
}

func (p *NimbusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config nimbusProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := defaultEndpoint
	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() && config.Endpoint.ValueString() != "" {
		endpoint = config.Endpoint.ValueString()
	}

	c := client.New(endpoint)
	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *NimbusProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *NimbusProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworkResource,
		NewSubnetResource,
		NewNetworkInterfaceResource,
		NewVMInstanceResource,
		NewStorageBucketResource,
	}
}
