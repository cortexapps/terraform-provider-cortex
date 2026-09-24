package integrations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}
var _ resource.ResourceWithConfigValidators = &Resource{}
var _ resource.ResourceWithValidateConfig = &Resource{}
var _ resource.ResourceWithModifyPlan = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// Resource manages one integration configuration. The integration definitions hold everything specific to an
// integration; this file holds the shared lifecycle.
type Resource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"alias": schema.StringAttribute{
			MarkdownDescription: "Unique alias of the configuration. Required for integrations that support several " +
				"configurations. A change renames the configuration in place.",
			Optional:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"is_default": schema.BoolAttribute{
			MarkdownDescription: "Whether this is the default configuration of its integration. When not set, Terraform " +
				"keeps the value from Cortex. Cortex makes the first configuration the default, and does not allow " +
				"setting the current default to `false`: set `is_default = true` on another configuration and remove " +
				"`is_default` from this one, apply, and then set it to `false` if necessary. Set `is_default = true` on " +
				"one configuration per integration only.",
			Optional: true,
			Computed: true,
			// No UseStateForUnknown: another configuration can take the default between plan and apply, so an
			// update reads the live value when is_default is not set.
			PlanModifiers: []planmodifier.Bool{defaultConfigurationModifier{}},
		},
		"credentials": credentialsAttribute(),
		"credentials_last_four": schema.MapAttribute{
			MarkdownDescription: "Last four characters of each secret credential part that Cortex stores, keyed by part " +
				"name (`value`, `password`, `key`, `secret`). Terraform uses them to detect a secret changed outside Terraform.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "`<integration>/<alias>`. Same as the import ID.",
			Computed:            true,
		},
		"integration": schema.StringAttribute{
			MarkdownDescription: "Name of the integration, from the settings block that is set.",
			Computed:            true,
		},
	}
	for _, d := range integrationDefinitions {
		attributes[d.Name()] = d.SettingsAttribute()
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Integration configuration. Set exactly one settings block to pick the integration " +
			"(`" + strings.ReplaceAll(definitionNames(), ", ", "`, `") + "`), and one credential kind in `credentials`.\n\n" +
			"| Integration | Credential kind | Mapping to the Cortex API |\n" +
			"|---|---|---|\n" +
			"| `datadog` | `key_pair` | `key` = API key, `secret` = application key |\n\n" +
			"The Cortex API never returns secrets. Terraform detects a secret changed outside Terraform through " +
			"`credentials_last_four`. Cortex does not check credentials when it saves a configuration, so invalid " +
			"credentials do not fail the apply. Cortex does not allow deleting the default configuration of an " +
			"integration while other configurations exist.",
		Attributes: attributes,
	}
}

func credentialsAttribute() schema.SingleNestedAttribute {
	secret := func(description string) schema.StringAttribute {
		return schema.StringAttribute{
			MarkdownDescription: description,
			Required:            true,
			Sensitive:           true,
			Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
		}
	}
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Credentials. Set exactly one kind. The integration defines which kind it accepts.",
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"token": schema.SingleNestedAttribute{
				MarkdownDescription: "One secret token or key.",
				Optional:            true,
				Attributes:          map[string]schema.Attribute{"value": secret("The token or key.")},
				// The validator counts the attribute it is on, so it sits on one kind and names the other two.
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("basic"),
						path.MatchRelative().AtParent().AtName("key_pair"),
					),
				},
			},
			"basic": schema.SingleNestedAttribute{
				MarkdownDescription: "A username or email, and a secret.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						MarkdownDescription: "The username or email.",
						Required:            true,
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"password": secret("The password or API token."),
				},
			},
			"key_pair": schema.SingleNestedAttribute{
				MarkdownDescription: "Two key parts that work together.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"key":    secret("The first key part."),
					"secret": secret("The second key part."),
				},
			},
		},
	}
}

// defaultConfigurationModifier fails the plan when the configuration sets is_default to false on the current
// default configuration, because Cortex rejects that update.
type defaultConfigurationModifier struct{}

func (m defaultConfigurationModifier) Description(ctx context.Context) string {
	return "Rejects is_default = false on the current default configuration."
}

func (m defaultConfigurationModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m defaultConfigurationModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if unsetsDefault(req.StateValue, req.ConfigValue) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Cannot unset the default configuration",
			"Cortex does not allow setting the current default configuration to is_default = false. Set is_default = true "+
				"on another configuration and remove is_default from this one, apply, and then set it to false if necessary.",
		)
	}
}

/***********************************************************************************************************************
 * Validation and plan
 **********************************************************************************************************************/

func (r *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := make([]path.Expression, 0, len(integrationDefinitions))
	for _, d := range integrationDefinitions {
		expressions = append(expressions, path.MatchRoot(d.Name()))
	}
	return []resource.ConfigValidator{resourcevalidator.ExactlyOneOf(expressions...)}
}

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	def := configuredDefinition(ctx, req.Config)
	if def == nil {
		return
	}
	var alias types.String
	var credentials types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("alias"), &alias)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("credentials"), &credentials)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isMultiInstance(def) && alias.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("alias"), "Missing alias",
			fmt.Sprintf("%s supports several configurations, so alias is required.", def.Title()))
	}

	if credentials.IsNull() || credentials.IsUnknown() {
		return
	}
	for name, value := range credentials.Attributes() {
		if !value.IsNull() && !allowsKind(def, credentialKind(name)) {
			resp.Diagnostics.AddAttributeError(path.Root("credentials").AtName(name), "Unsupported credential kind",
				fmt.Sprintf("%s accepts: %s.", def.Name(), kindNames(def.CredentialKinds())))
		}
	}
}

// ModifyPlan sets the computed values that the plan can know, and marks the changes that need a new configuration.
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan integrationConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	def := configuredDefinition(ctx, req.Plan)
	if def == nil {
		return
	}

	plan.Integration = types.StringValue(def.Name())
	if !plan.Alias.IsNull() && !plan.Alias.IsUnknown() {
		plan.Id = types.StringValue(def.Name() + "/" + plan.Alias.ValueString())
	}

	if req.State.Raw.IsNull() {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		return
	}
	var state integrationConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Integration.ValueString() != def.Name() {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("integration"))
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		return
	}

	credentialsChanged := false
	if state.Credentials != nil && plan.Credentials != nil && state.Credentials.kind() != plan.Credentials.kind() {
		credentialsChanged = true
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("credentials"))
	} else if plan.Credentials != nil {
		stateLastFour, diags := stringMap(ctx, state.CredentialsLastFour)
		resp.Diagnostics.Append(diags...)
		for _, p := range credentialParts[plan.Credentials.kind()] {
			if !p.secret {
				continue
			}
			statePart := types.StringNull()
			if s := state.Credentials.part(p.name); s != nil {
				statePart = *s
			}
			want, ok := stateLastFour[p.name]
			if secretChanged(statePart, *plan.Credentials.part(p.name), want, ok) {
				credentialsChanged = true
			}
		}
		if credentialsChanged && !def.CredentialsUpdatable() {
			resp.RequiresReplace = append(resp.RequiresReplace, path.Root("credentials"))
		}
	}
	if !credentialsChanged {
		plan.CredentialsLastFour = state.CredentialsLastFour
	}

	replace, diags := def.SettingsRequireReplace(ctx, *plan.settings(def.Name()), *state.settings(def.Name()))
	resp.Diagnostics.Append(diags...)
	if replace {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root(def.Name()))
	}
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_configuration"
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*cortex.HttpClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	def := definitionByName(state.Integration.ValueString())
	if def == nil {
		resp.Diagnostics.AddError("Unknown integration", fmt.Sprintf("State has no known integration. Use one of: %s.", definitionNames()))
		return
	}
	prior := *state.settings(def.Name())

	var st configurationState
	switch d := def.(type) {
	case multiInstanceDefinition:
		states, err := d.List(ctx, r.client, prior)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s configuration %s, got error: %s", def.Name(), state.Alias.ValueString(), err))
			return
		}
		found, ok := findState(states, state.Alias.ValueString())
		if !ok {
			resp.State.RemoveResource(ctx)
			return
		}
		st = found
	default:
		addUnsupportedEngineError(&resp.Diagnostics, def)
		return
	}

	applyState(ctx, &state, def, st, &resp.Diagnostics)
	clearDriftedSecrets(state.Credentials, st.LastFour)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	def := configuredDefinition(ctx, req.Plan)
	in := configurationInput{Alias: plan.Alias.ValueString(), Credentials: plan.Credentials.value(), Settings: *plan.settings(def.Name())}

	var st configurationState
	var err error
	switch d := def.(type) {
	case multiInstanceDefinition:
		// Cortex makes a new configuration the default when no default exists, whatever the request asks. Fail
		// before the create, so an explicit is_default = false does not leave a configuration that contradicts the plan.
		if knownBoolEquals(plan.IsDefault, false) {
			states, lerr := d.List(ctx, r.client, in.Settings)
			if lerr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list %s configurations, got error: %s", def.Name(), lerr))
				return
			}
			if noDefault(states) {
				resp.Diagnostics.AddAttributeError(path.Root("is_default"), "Configuration must be the default",
					fmt.Sprintf("No %s configuration is the default, so Cortex makes %s the default. Remove is_default or set it to true.", def.Name(), in.Alias))
				return
			}
		}
		// The create request never asks for the default: Cortex rejects a second default. An update sets it.
		st, err = d.Create(ctx, r.client, in)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s configuration, got error: %s", def.Name(), err))
			return
		}
		if knownBoolEquals(plan.IsDefault, true) && !st.IsDefault {
			in.IsDefault = true
			updated, uerr := d.Update(ctx, r.client, st.Alias, in)
			if uerr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Created %s configuration %s, but unable to make it the default, got error: %s", def.Name(), st.Alias, uerr))
			} else {
				st = updated
			}
		}
	default:
		addUnsupportedEngineError(&resp.Diagnostics, def)
		return
	}

	// Save state also after an error, so Terraform keeps track of the configuration.
	applyState(ctx, &plan, def, st, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state integrationConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	def := configuredDefinition(ctx, req.Plan)
	in := configurationInput{Alias: plan.Alias.ValueString(), IsDefault: plan.IsDefault.ValueBool(), Credentials: plan.Credentials.value(), Settings: *plan.settings(def.Name())}

	var st configurationState
	var err error
	switch d := def.(type) {
	case multiInstanceDefinition:
		// When is_default is not set, send the live value, so the update does not take back a default that another
		// configuration took after the plan.
		if plan.IsDefault.IsUnknown() {
			states, lerr := d.List(ctx, r.client, in.Settings)
			if lerr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s configuration %s, got error: %s", def.Name(), state.Alias.ValueString(), lerr))
				return
			}
			live, ok := findState(states, state.Alias.ValueString())
			if !ok {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s configuration %s does not exist anymore.", def.Title(), state.Alias.ValueString()))
				return
			}
			in.IsDefault = live.IsDefault
		}
		st, err = d.Update(ctx, r.client, state.Alias.ValueString(), in)
	default:
		addUnsupportedEngineError(&resp.Diagnostics, def)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update %s configuration, got error: %s", def.Name(), err))
		return
	}

	applyState(ctx, &plan, def, st, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	def := definitionByName(state.Integration.ValueString())
	var err error
	switch d := def.(type) {
	case multiInstanceDefinition:
		err = d.Delete(ctx, r.client, state.Alias.ValueString())
	default:
		addUnsupportedEngineError(&resp.Diagnostics, def)
		return
	}
	if err != nil && !errors.Is(err, cortex.ApiErrorNotFound) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s configuration, got error: %s", def.Name(), err))
	}
}

// ImportState takes <integration>/<alias>.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	name, alias, hasAlias := strings.Cut(req.ID, "/")
	def := definitionByName(name)
	if def == nil {
		resp.Diagnostics.AddError("Unknown integration", fmt.Sprintf("Unknown integration %q in import ID. Use one of: %s.", name, definitionNames()))
		return
	}
	if isMultiInstance(def) && (!hasAlias || alias == "") {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("The import ID for %s is %s/<alias>.", def.Title(), def.Name()))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("integration"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	if hasAlias {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("alias"), alias)...)
	}
}

/***********************************************************************************************************************
 * State
 **********************************************************************************************************************/

// applyState copies what the API returned into the model.
func applyState(ctx context.Context, m *integrationConfigurationModel, def integrationDefinition, st configurationState, diags *diag.Diagnostics) {
	m.Integration = types.StringValue(def.Name())
	m.Alias = types.StringValue(st.Alias)
	m.IsDefault = types.BoolValue(st.IsDefault)
	m.Id = types.StringValue(def.Name() + "/" + st.Alias)
	*m.settings(def.Name()) = st.Settings

	lastFours := st.LastFour
	if lastFours == nil {
		lastFours = map[string]string{}
	}
	value, d := types.MapValueFrom(ctx, types.StringType, lastFours)
	diags.Append(d...)
	m.CredentialsLastFour = value
}

// addUnsupportedEngineError reports a definition that no engine handles, instead of writing an empty state.
func addUnsupportedEngineError(diags *diag.Diagnostics, def integrationDefinition) {
	name := "unknown"
	if def != nil {
		name = def.Name()
	}
	diags.AddError("Unsupported integration", fmt.Sprintf("The provider has no engine for integration %q. Please report this issue to the provider developers.", name))
}

func noDefault(states []configurationState) bool {
	for _, s := range states {
		if s.IsDefault {
			return false
		}
	}
	return true
}
