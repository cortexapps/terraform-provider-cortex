package integrations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

/***********************************************************************************************************************
 * Credentials
 **********************************************************************************************************************/

type credentialKind string

const (
	credentialToken   credentialKind = "token"
	credentialBasic   credentialKind = "basic"
	credentialKeyPair credentialKind = "key_pair"
)

// credentialPart is one value of a credential kind. The API never returns a secret part.
type credentialPart struct {
	name   string
	secret bool
}

var credentialParts = map[credentialKind][]credentialPart{
	credentialToken:   {{name: "value", secret: true}},
	credentialBasic:   {{name: "username", secret: false}, {name: "password", secret: true}},
	credentialKeyPair: {{name: "key", secret: true}, {name: "secret", secret: true}},
}

// credentialsModel is the decoded credentials object. The resource model keeps the credentials as a types.Object,
// because a Go pointer cannot hold a value that is unknown until apply. decodeCredentials builds this model once the
// values are known, and object converts it back.
type credentialsModel struct {
	Token   *tokenCredentialModel   `tfsdk:"token"`
	Basic   *basicCredentialModel   `tfsdk:"basic"`
	KeyPair *keyPairCredentialModel `tfsdk:"key_pair"`
}

// credentialsObjectModel is the credentials object with one types.Object per kind, which can be null or unknown.
type credentialsObjectModel struct {
	Token   types.Object `tfsdk:"token"`
	Basic   types.Object `tfsdk:"basic"`
	KeyPair types.Object `tfsdk:"key_pair"`
}

var credentialKindAttrTypes = map[credentialKind]map[string]attr.Type{
	credentialToken:   {"value": types.StringType},
	credentialBasic:   {"username": types.StringType, "password": types.StringType},
	credentialKeyPair: {"key": types.StringType, "secret": types.StringType},
}

var credentialsAttrTypes = map[string]attr.Type{
	string(credentialToken):   types.ObjectType{AttrTypes: credentialKindAttrTypes[credentialToken]},
	string(credentialBasic):   types.ObjectType{AttrTypes: credentialKindAttrTypes[credentialBasic]},
	string(credentialKeyPair): types.ObjectType{AttrTypes: credentialKindAttrTypes[credentialKeyPair]},
}

// decodeCredentials converts the credentials object to the model. It returns a nil model when the object is null.
// known is false, with a nil model, when the object or the object of a kind is unknown. Parts inside a known kind
// object can still be unknown; the model keeps them as unknown strings.
func decodeCredentials(ctx context.Context, obj types.Object) (c *credentialsModel, known bool, diags diag.Diagnostics) {
	if obj.IsNull() {
		return nil, true, diags
	}
	if obj.IsUnknown() {
		return nil, false, diags
	}
	var o credentialsObjectModel
	diags.Append(obj.As(ctx, &o, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, true, diags
	}
	for _, kind := range []types.Object{o.Token, o.Basic, o.KeyPair} {
		if kind.IsUnknown() {
			return nil, false, diags
		}
	}
	c = &credentialsModel{}
	if !o.Token.IsNull() {
		c.Token = &tokenCredentialModel{}
		diags.Append(o.Token.As(ctx, c.Token, basetypes.ObjectAsOptions{})...)
	}
	if !o.Basic.IsNull() {
		c.Basic = &basicCredentialModel{}
		diags.Append(o.Basic.As(ctx, c.Basic, basetypes.ObjectAsOptions{})...)
	}
	if !o.KeyPair.IsNull() {
		c.KeyPair = &keyPairCredentialModel{}
		diags.Append(o.KeyPair.As(ctx, c.KeyPair, basetypes.ObjectAsOptions{})...)
	}
	return c, true, diags
}

// object converts the model back to the credentials object. A nil model gives a null object.
func (c *credentialsModel) object(ctx context.Context) (types.Object, diag.Diagnostics) {
	if c == nil {
		return types.ObjectNull(credentialsAttrTypes), nil
	}
	var diags diag.Diagnostics
	kind := func(k credentialKind, set bool, value any) types.Object {
		t := credentialKindAttrTypes[k]
		if !set {
			return types.ObjectNull(t)
		}
		o, d := types.ObjectValueFrom(ctx, t, value)
		diags.Append(d...)
		return o
	}
	o, d := types.ObjectValueFrom(ctx, credentialsAttrTypes, credentialsObjectModel{
		Token:   kind(credentialToken, c.Token != nil, c.Token),
		Basic:   kind(credentialBasic, c.Basic != nil, c.Basic),
		KeyPair: kind(credentialKeyPair, c.KeyPair != nil, c.KeyPair),
	})
	diags.Append(d...)
	return o, diags
}

type tokenCredentialModel struct {
	Value types.String `tfsdk:"value"`
}

type basicCredentialModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type keyPairCredentialModel struct {
	Key    types.String `tfsdk:"key"`
	Secret types.String `tfsdk:"secret"`
}

// credentialsValue holds the credentials as plain values for the integration definitions.
type credentialsValue struct {
	Kind  credentialKind
	Parts map[string]string
}

func (c *credentialsModel) kind() credentialKind {
	switch {
	case c == nil:
		return ""
	case c.Token != nil:
		return credentialToken
	case c.Basic != nil:
		return credentialBasic
	case c.KeyPair != nil:
		return credentialKeyPair
	}
	return ""
}

// part returns the value of a part of the set kind, or nil when the kind has no such part.
func (c *credentialsModel) part(name string) *types.String {
	switch c.kind() {
	case credentialToken:
		if name == "value" {
			return &c.Token.Value
		}
	case credentialBasic:
		switch name {
		case "username":
			return &c.Basic.Username
		case "password":
			return &c.Basic.Password
		}
	case credentialKeyPair:
		switch name {
		case "key":
			return &c.KeyPair.Key
		case "secret":
			return &c.KeyPair.Secret
		}
	}
	return nil
}

func (c *credentialsModel) value() credentialsValue {
	v := credentialsValue{Kind: c.kind(), Parts: map[string]string{}}
	for _, p := range credentialParts[v.Kind] {
		if s := c.part(p.name); s != nil {
			v.Parts[p.name] = s.ValueString()
		}
	}
	return v
}

/***********************************************************************************************************************
 * Secret drift
 **********************************************************************************************************************/

func lastFour(value string) string {
	runes := []rune(value)
	if len(runes) <= 4 {
		return value
	}
	return string(runes[len(runes)-4:])
}

// secretChanged tells if a secret in the plan differs from the secret that Cortex stores. It compares with the
// value in state, and with the last four characters from the API. The second check catches a secret rotated outside
// Terraform, and a configured secret that does not match Cortex after an import (the state value is then null).
func secretChanged(stateSecret types.String, planSecret types.String, stateLastFour string, hasLastFour bool) bool {
	if planSecret.IsUnknown() {
		return true
	}
	if !stateSecret.IsNull() && !stateSecret.IsUnknown() && stateSecret.ValueString() != planSecret.ValueString() {
		return true
	}
	return hasLastFour && lastFour(planSecret.ValueString()) != stateLastFour
}

// clearDriftedSecrets removes each secret whose last four characters no longer match the API value. Terraform
// replaces or updates a resource only when a value changes, so the null value makes the next plan show the change.
// Call it only on Read.
func clearDriftedSecrets(c *credentialsModel, apiLastFour map[string]string) {
	for _, p := range credentialParts[c.kind()] {
		if !p.secret {
			continue
		}
		s := c.part(p.name)
		want, ok := apiLastFour[p.name]
		if s == nil || !ok || s.IsNull() || s.IsUnknown() {
			continue
		}
		if lastFour(s.ValueString()) != want {
			*s = types.StringNull()
		}
	}
}

/***********************************************************************************************************************
 * Default configuration
 **********************************************************************************************************************/

// unsetsDefault tells if the configuration sets is_default to false on the current default configuration. Cortex
// rejects this: a different configuration must become the default first.
func unsetsDefault(stateValue types.Bool, configValue types.Bool) bool {
	if stateValue.IsNull() || stateValue.IsUnknown() || configValue.IsNull() || configValue.IsUnknown() {
		return false
	}
	return stateValue.ValueBool() && !configValue.ValueBool()
}

// knownBoolEquals tells if the value is known, not null, and equal to want. It separates an explicit configuration
// value from an omitted one.
func knownBoolEquals(value types.Bool, want bool) bool {
	return !value.IsNull() && !value.IsUnknown() && value.ValueBool() == want
}

/***********************************************************************************************************************
 * Value helpers
 **********************************************************************************************************************/

func stringMap(ctx context.Context, m types.Map) (map[string]string, diag.Diagnostics) {
	out := map[string]string{}
	if m.IsNull() || m.IsUnknown() {
		return out, nil
	}
	diags := m.ElementsAs(ctx, &out, false)
	return out, diags
}

// stringList never returns nil, because the API rejects a null list.
func stringList(ctx context.Context, l types.List) ([]string, diag.Diagnostics) {
	out := []string{}
	if l.IsNull() || l.IsUnknown() {
		return out, nil
	}
	diags := l.ElementsAs(ctx, &out, false)
	return out, diags
}

func optionalString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

/***********************************************************************************************************************
 * Resource model
 **********************************************************************************************************************/

type integrationConfigurationModel struct {
	Id                  types.String `tfsdk:"id"`
	Integration         types.String `tfsdk:"integration"`
	Alias               types.String `tfsdk:"alias"`
	IsDefault           types.Bool   `tfsdk:"is_default"`
	Credentials         types.Object `tfsdk:"credentials"`
	CredentialsLastFour types.Map    `tfsdk:"credentials_last_four"`
	Datadog             types.Object `tfsdk:"datadog"`
	Gitlab              types.Object `tfsdk:"gitlab"`
	IncidentIo          types.Object `tfsdk:"incident_io"`
	PagerDuty           types.Object `tfsdk:"pagerduty"`
}

// settings returns the settings block of the named integration. Every registered definition needs a case here.
func (m *integrationConfigurationModel) settings(name string) *types.Object {
	switch name {
	case "datadog":
		return &m.Datadog
	case "gitlab":
		return &m.Gitlab
	case "incident_io":
		return &m.IncidentIo
	case "pagerduty":
		return &m.PagerDuty
	}
	return nil
}
