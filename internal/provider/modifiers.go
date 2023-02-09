package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type stringDefaultModifier struct {
	Default string
}

func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

func (m stringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	ctx = tflog.SetField(ctx, "__plan_null", req.PlanValue.IsNull())
	ctx = tflog.SetField(ctx, "__plan_unknown", req.PlanValue.IsUnknown())
	ctx = tflog.SetField(ctx, "__state_null", req.StateValue.IsNull())
	ctx = tflog.SetField(ctx, "__state_unknown", req.StateValue.IsUnknown())
	ctx = tflog.SetField(ctx, "__config_null", req.ConfigValue.IsNull())
	ctx = tflog.SetField(ctx, "__config_unknown", req.ConfigValue.IsUnknown())
	ctx = tflog.SetField(ctx, "__default_value", m.Default)

	if !req.PlanValue.IsNull() && !(req.ConfigValue.IsNull() && req.PlanValue.IsUnknown()) {
		tflog.Info(ctx, "Not setting default value")
		return
	}

	resp.PlanValue = types.StringValue(m.Default)
	tflog.Info(ctx, "Set default value")
}

type int64DefaultModifier struct {
	Default int64
}

func (m int64DefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %d", m.Default)
}

func (m int64DefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%d`", m.Default)
}

func (m int64DefaultModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.PlanValue.IsNull() && !(req.ConfigValue.IsNull() && req.PlanValue.IsUnknown()) {
		return
	}

	resp.PlanValue = types.Int64Value(m.Default)
}

type boolDefaultModifier struct {
	Default bool
}

func (m boolDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %t", m.Default)
}

func (m boolDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%t`", m.Default)
}

func (m boolDefaultModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	ctx = tflog.SetField(ctx, "__plan_null", req.PlanValue.IsNull())
	ctx = tflog.SetField(ctx, "__plan_unknown", req.PlanValue.IsUnknown())
	ctx = tflog.SetField(ctx, "__state_null", req.StateValue.IsNull())
	ctx = tflog.SetField(ctx, "__state_unknown", req.StateValue.IsUnknown())
	ctx = tflog.SetField(ctx, "__config_null", req.ConfigValue.IsNull())
	ctx = tflog.SetField(ctx, "__config_unknown", req.ConfigValue.IsUnknown())
	ctx = tflog.SetField(ctx, "__default_value", m.Default)

	if !req.PlanValue.IsNull() && !(req.ConfigValue.IsNull() && req.PlanValue.IsUnknown()) {
		tflog.Info(ctx, "Not setting default value")
		return
	}

	resp.PlanValue = types.BoolValue(m.Default)
	tflog.Info(ctx, "Set default value")
}

type testModifier struct {
	Default string
}

func (m testModifier) Description(ctx context.Context) string {
	return "logging"
}

func (m testModifier) MarkdownDescription(ctx context.Context) string {
	return "logging"
}

func (m testModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {

	var state, plan types.Bool
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("generate_s3_credentials"), &state)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("generate_s3_credentials"), &plan)...)

	ctx = tflog.SetField(ctx, "__generate_credentials_plan", plan.ValueBool())
	ctx = tflog.SetField(ctx, "__generate_credentials_plan_null", plan.IsNull())
	ctx = tflog.SetField(ctx, "__generate_credentials_plan_unknown", plan.IsUnknown())
	ctx = tflog.SetField(ctx, "__generate_credentials_state", state.ValueBool())
	ctx = tflog.SetField(ctx, "__generate_credentials_state_null", state.IsNull())
	ctx = tflog.SetField(ctx, "__generate_credentials_state_unknown", state.IsUnknown())
	tflog.Info(ctx, "no doing anything")
}
