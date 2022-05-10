package vercel

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vercel/terraform-provider-vercel/client"
)

// Deployment represents the terraform state for a deployment resource.
type Alias struct {
	Alias        types.String `tfsdk:"alias"`
	DeploymentId types.String `tfsdk:"deployment_id"`
	TeamId       types.String `tfsdk:"team_id"`
	Production   types.Bool   `tfsdk:"production"`
}

// convertResponseToDeployment is used to populate terraform state based on an API response.
// Where possible, values from the API response are used to populate state. If not possible,
// values from the existing deployment state are used.
func convertResponseToAlias(response client.AliasResponse, plan Alias) Alias {
	production := types.Bool{Value: false}
	/*
	 * TODO - the first deployment to a new project is currently _always_ a
	 * production deployment, even if you ask it to be a preview deployment.
	 * In order to terraform complaining about an inconsistent output, we should only set
	 * the state back if it matches what we expect. The third part of this
	 * conditional ensures this, but can be removed if the behaviour is changed.
	 * see:
	 * https://github.com/vercel/customer-issues/issues/178#issuecomment-1012062345 and
	 * https://vercel.slack.com/archives/C01A2M9R8RZ/p1639594164360300
	 * for more context.
	 */
	if response.Target != nil && *response.Target == "production" && (plan.Production.Value || plan.Production.Unknown) {
		production.Value = true
	}

	return Alias{
		Alias:        types.String{Value: response.Alias},
		DeploymentId: plan.DeploymentId,
		TeamID:       plan.TeamId,
	}
}
