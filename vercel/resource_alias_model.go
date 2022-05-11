package vercel

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vercel/terraform-provider-vercel/client"
)

// Alias represents the terraform state for an alias resource.
type Alias struct {
	Alias        types.String `tfsdk:"alias"`
	AliasUID     types.String `tfsdk:"-"`
	DeploymentId types.String `tfsdk:"deployment_id"`
	TeamID       types.String `tfsdk:"team_id"`
}

// convertResponseToDeployment is used to populate terraform state based on an API response.
// Where possible, values from the API response are used to populate state. If not possible,
// values from the existing deployment state are used.
func convertResponseToAlias(response client.AliasResponse, plan Alias) Alias {
	return Alias{
		Alias:        types.String{Value: response.Alias},
		AliasUID:     types.String{Value: response.UID},
		DeploymentId: plan.DeploymentId,
		TeamID:       plan.TeamID,
	}
}
