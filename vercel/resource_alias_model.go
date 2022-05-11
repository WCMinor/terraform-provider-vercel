package vercel

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vercel/terraform-provider-vercel/client"
)

// Deployment represents the terraform state for a deployment resource.
type Alias struct {
	Alias        types.String `tfsdk:"alias"`
	DeploymentId types.String `tfsdk:"deployment_id"`
	TeamID       types.String `tfsdk:"team_id"`
	Production   types.Bool   `tfsdk:"production"`
	AliasUID     types.String `tfsdk:"-"`
}

// convertResponseToDeployment is used to populate terraform state based on an API response.
// Where possible, values from the API response are used to populate state. If not possible,
// values from the existing deployment state are used.
func convertResponseToAlias(response client.CreateAliasResponse, plan Alias) Alias {
	return Alias{
		Alias:        types.String{Value: response.Alias},
		DeploymentId: plan.DeploymentId,
		AliasUID:     types.String{Value: response.UID},
		TeamID:       plan.TeamID,
	}
}
