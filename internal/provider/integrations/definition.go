package integrations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

/***********************************************************************************************************************
 * Definitions
 **********************************************************************************************************************/

// configurationInput is what the resource sends to a definition.
type configurationInput struct {
	Alias       string // multi-instance only
	IsDefault   bool   // multi-instance only
	Credentials credentialsValue
	Settings    types.Object // the settings block from the plan
}

// configurationState is what a definition reads back from the API, in provider terms.
type configurationState struct {
	Alias     string            // multi-instance only
	IsDefault bool              // multi-instance only
	Settings  types.Object      // the settings block
	LastFour  map[string]string // secret credential part -> last four characters
}

// integrationDefinition describes one integration. To add an integration, write a definition, add it to
// integrationDefinitions, and add its settings field to integrationConfigurationModel.
type integrationDefinition interface {
	// Name is the settings block name and the import prefix, for example "datadog".
	Name() string
	// Title is the name in messages, for example "Datadog".
	Title() string
	CredentialKinds() []credentialKind
	SettingsAttribute() schema.SingleNestedAttribute
	// SettingsRequireReplace tells if a settings change needs a new configuration. When settingsUnknown is true, it
	// must return true if the unknown block can change a field that needs a new configuration.
	SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics)
}

// multiInstanceDefinition is an integration with several configurations per tenant, identified by alias.
type multiInstanceDefinition interface {
	integrationDefinition
	// List returns the configurations the provider supports. prior is the settings block from state, for values
	// the API does not return.
	List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error)
	Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error)
	Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error)
	Delete(ctx context.Context, c *cortex.HttpClient, alias string) error
}

// singleInstanceDefinition is an integration with at most one configuration per tenant.
type singleInstanceDefinition interface {
	integrationDefinition
	// Get returns an error that matches cortex.ApiErrorNotFound when no configuration exists.
	Get(ctx context.Context, c *cortex.HttpClient, prior types.Object) (configurationState, error)
	Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error)
	// Replace sends the full configuration.
	Replace(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error)
	Delete(ctx context.Context, c *cortex.HttpClient) error
}

var integrationDefinitions = []integrationDefinition{
	datadogDefinition{},
	gitlabDefinition{},
	incidentIoDefinition{},
	pagerDutyDefinition{},
}

func definitionByName(name string) integrationDefinition {
	for _, d := range integrationDefinitions {
		if d.Name() == name {
			return d
		}
	}
	return nil
}

func definitionNames() string {
	names := make([]string, 0, len(integrationDefinitions))
	for _, d := range integrationDefinitions {
		names = append(names, d.Name())
	}
	return strings.Join(names, ", ")
}

func isMultiInstance(d integrationDefinition) bool {
	_, ok := d.(multiInstanceDefinition)
	return ok
}

func isSingleInstance(d integrationDefinition) bool {
	_, ok := d.(singleInstanceDefinition)
	return ok
}

func allowsKind(d integrationDefinition, kind credentialKind) bool {
	for _, k := range d.CredentialKinds() {
		if k == kind {
			return true
		}
	}
	return false
}

func kindNames(kinds []credentialKind) string {
	names := make([]string, 0, len(kinds))
	for _, k := range kinds {
		names = append(names, string(k))
	}
	return strings.Join(names, ", ")
}

func findState(states []configurationState, alias string) (configurationState, bool) {
	for _, s := range states {
		if s.Alias == alias {
			return s, true
		}
	}
	return configurationState{}, false
}

// attributeGetter is the part of tfsdk.Config, tfsdk.Plan, and tfsdk.State that the resource uses.
type attributeGetter interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

// configuredDefinition returns the definition whose settings block is set, or nil.
func configuredDefinition(ctx context.Context, src attributeGetter) integrationDefinition {
	for _, d := range integrationDefinitions {
		var v types.Object
		if diags := src.GetAttribute(ctx, path.Root(d.Name()), &v); diags.HasError() {
			continue
		}
		if !v.IsNull() {
			return d
		}
	}
	return nil
}

/***********************************************************************************************************************
 * Settings helpers
 **********************************************************************************************************************/

func settingsAttrTypes(d integrationDefinition) map[string]attr.Type {
	return d.SettingsAttribute().GetType().(types.ObjectType).AttrTypes
}

// settingsFrom converts a settings block into the definition's settings struct.
func settingsFrom[T any](ctx context.Context, obj types.Object) (T, error) {
	var v T
	diags := obj.As(ctx, &v, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	return v, diagsError(diags)
}

// settingsObject converts the definition's settings struct into a settings block.
func settingsObject[T any](ctx context.Context, d integrationDefinition, v T) (types.Object, error) {
	obj, diags := types.ObjectValueFrom(ctx, settingsAttrTypes(d), v)
	return obj, diagsError(diags)
}

// settingsUnknown tells if the planned settings block is unknown for an existing configuration, for example because
// it comes from another resource. The apply can then find a change to any field, so a replace check must assume one.
func settingsUnknown(plan types.Object, state types.Object) bool {
	return plan.IsUnknown() && !state.IsNull() && !state.IsUnknown()
}

// asSettings converts the plan and state settings blocks for a replace check. It returns both values or neither:
// both are nil when either block is null or unknown, so a caller only needs to check one. Check settingsUnknown
// first: a nil result does not mean that nothing changes.
func asSettings[T any](ctx context.Context, plan types.Object, state types.Object) (*T, *T, diag.Diagnostics) {
	var diags diag.Diagnostics
	if plan.IsNull() || plan.IsUnknown() || state.IsNull() || state.IsUnknown() {
		return nil, nil, diags
	}
	var p, s T
	diags.Append(plan.As(ctx, &p, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)
	diags.Append(state.As(ctx, &s, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, nil, diags
	}
	return &p, &s, diags
}

func diagsError(diags diag.Diagnostics) error {
	if !diags.HasError() {
		return nil
	}
	messages := []string{}
	for _, d := range diags.Errors() {
		messages = append(messages, d.Summary()+": "+d.Detail())
	}
	return errors.New(strings.Join(messages, "; "))
}

func settingsError(d integrationDefinition, err error) error {
	return fmt.Errorf("invalid %s settings: %w", d.Name(), err)
}
