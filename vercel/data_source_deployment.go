package vercel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceDeploymentType struct{}

// GetSchema returns the schema information for a deployment data source
func (r dataSourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides information about an existing deployment within Vercel.

A Deployment is the result of building your Project and making it available through a live URL.
        `,
		Attributes: map[string]tfsdk.Attribute{
			"domains": {
				Description: "A list of all the domains (default domains, staging domains and production domains) that were assigned upon deployment creation.",
				Computed:    true,
				Type: types.ListType{
					ElemType: types.StringType,
				},
			},
			"environment": {
				Description: "A map of environment variable names to values. These are specific to a Deployment, and can also be configured on the `vercel_project` resource.",
				Computed:    true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
			},
			"team_id": {
				Description: "The team ID to add the deployment to.",
				Optional:    true,
				Type:        types.StringType,
			},
			"project_id": {
				Description: "The project ID to add the deployment to.",
				Required:    true,
				Type:        types.StringType,
			},
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"path_prefix": {
				Description: "File paths as they are uploaded to Vercel.",
				Computed:    true,
				Type:        types.StringType,
			},
			"url": {
				Description: "A unique URL that is automatically generated for a deployment.",
				Computed:    true,
				Type:        types.StringType,
			},
			"production": {
				Description: "true if the deployment is a production deployment.",
				Computed:    true,
				Type:        types.BoolType,
			},
			"files": {
				Description: "A map of files uploaded for the deployment.",
				Computed:    true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
			},
			"sha": {
				Description: "The specific commit hash that was used for the deployment. Note this will only work if the project is configured to use a Git repository.",
				Computed:    true,
				Type:        types.StringType,
			},
			"ref": {
				Description: "The branch or commit hash used for the deployment. Note this will only work if the project is configured to use a Git repository.",
				Computed:    true,
				Type:        types.StringType,
			},
			"project_settings": {
				Description: "Project settings applied to the deployment.",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"build_command": {
						Computed:    true,
						Type:        types.StringType,
						Description: "The build command for this deployment.",
					},
					"framework": {
						Computed:    true,
						Type:        types.StringType,
						Description: "The framework that is being used for this deployment.",
					},
					"install_command": {
						Computed:    true,
						Type:        types.StringType,
						Description: "The install command for this deployment.",
					},
					"output_directory": {
						Computed:    true,
						Type:        types.StringType,
						Description: "The output directory of the deployment.",
					},
					"root_directory": {
						Computed:    true,
						Type:        types.StringType,
						Description: "The name of a directory or relative path to the source code of your project.",
					},
				}),
			},
			"delete_on_destroy": {
				Description: "Hard delete the Vercel deployment when destroying the Terraform resource.",
				Computed:    true,
				Type:        types.BoolType,
			},
		},
	}, nil
}

// NewDataSource instantiates a new DataSource of this DataSourceType.
func (r dataSourceDeploymentType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceDeployment{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceDeployment struct {
	p provider
}

// Read will read the deployment information by requesting it from the Vercel API, and will update terraform
// with this information.
// It is called by the provider whenever data source values should be read to update state.
func (r dataSourceDeployment) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config Deployment
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.p.client.GetDeployment(ctx, config.ID.Value, config.TeamID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading deployment",
			fmt.Sprintf("Could not read deployment %s %s, unexpected error: %s",
				config.TeamID.Value,
				config.ID.Value,
				err,
			),
		)
		return
	}

	result := convertResponseToDeployment(out, config)
	tflog.Trace(ctx, "read deployment", map[string]interface{}{
		"team_id":       result.TeamID.Value,
		"deployment_id": result.ID.Value,
	})

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
