package integrations

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialsModelKindAndParts(t *testing.T) {
	c := &credentialsModel{KeyPair: &keyPairCredentialModel{Key: types.StringValue("k-1234"), Secret: types.StringValue("s-5678")}}
	assert.Equal(t, credentialKeyPair, c.kind())
	assert.Equal(t, "k-1234", c.part("key").ValueString())
	assert.Nil(t, c.part("value"))
	assert.Equal(t, credentialsValue{Kind: credentialKeyPair, Parts: map[string]string{"key": "k-1234", "secret": "s-5678"}}, c.value())

	var none *credentialsModel
	assert.Equal(t, credentialKind(""), none.kind())
	assert.Nil(t, none.part("value"))
}

func TestSecretChanged(t *testing.T) {
	tests := []struct {
		name        string
		state, plan types.String
		lastFour    string
		hasLastFour bool
		want        bool
	}{
		{"unchanged", types.StringValue("old-1234"), types.StringValue("old-1234"), "1234", true, false},
		{"changed", types.StringValue("old-1234"), types.StringValue("new-9999"), "1234", true, true},
		{"changed, same last four", types.StringValue("old-1234"), types.StringValue("new-1234"), "1234", true, true},
		{"rotated outside terraform", types.StringValue("old-1234"), types.StringValue("old-1234"), "9999", true, true},
		{"imported, config matches api", types.StringNull(), types.StringValue("old-1234"), "1234", true, false},
		{"imported, config does not match api", types.StringNull(), types.StringValue("new-9999"), "1234", true, true},
		{"unknown plan", types.StringValue("old-1234"), types.StringUnknown(), "1234", true, true},
		{"no last four", types.StringValue("old-1234"), types.StringValue("old-1234"), "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, secretChanged(tt.state, tt.plan, tt.lastFour, tt.hasLastFour))
		})
	}
}

func TestClearDriftedSecrets(t *testing.T) {
	c := &credentialsModel{Basic: &basicCredentialModel{Username: types.StringValue("bot"), Password: types.StringValue("pass-1234")}}
	clearDriftedSecrets(c, map[string]string{"password": "9999"})
	assert.True(t, c.Basic.Password.IsNull())
	assert.Equal(t, "bot", c.Basic.Username.ValueString())

	kept := &credentialsModel{Token: &tokenCredentialModel{Value: types.StringValue("tok-1234")}}
	clearDriftedSecrets(kept, map[string]string{"value": "1234"})
	assert.Equal(t, "tok-1234", kept.Token.Value.ValueString())

	clearDriftedSecrets(nil, map[string]string{"value": "1234"})
}

func TestLastFour(t *testing.T) {
	assert.Equal(t, "1234", lastFour("abc-1234"))
	assert.Equal(t, "ab", lastFour("ab"))
}

func TestUnsetsDefault(t *testing.T) {
	assert.True(t, unsetsDefault(types.BoolValue(true), types.BoolValue(false)))
	assert.False(t, unsetsDefault(types.BoolValue(true), types.BoolNull()))
	assert.False(t, unsetsDefault(types.BoolValue(false), types.BoolValue(false)))
	assert.False(t, unsetsDefault(types.BoolNull(), types.BoolValue(false)))
}

func TestKnownBoolEquals(t *testing.T) {
	assert.True(t, knownBoolEquals(types.BoolValue(true), true))
	assert.False(t, knownBoolEquals(types.BoolNull(), false))
	assert.False(t, knownBoolEquals(types.BoolUnknown(), false))
}

func TestStringHelpers(t *testing.T) {
	ctx := context.Background()
	m, d := stringMap(ctx, types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("1234")}))
	require.False(t, d.HasError())
	assert.Equal(t, map[string]string{"key": "1234"}, m)
	empty, _ := stringMap(ctx, types.MapNull(types.StringType))
	assert.Equal(t, map[string]string{}, empty)

	l, d := stringList(ctx, types.ListNull(types.StringType))
	require.False(t, d.HasError())
	assert.NotNil(t, l)
	assert.Empty(t, l)

	assert.True(t, optionalString("").IsNull())
	assert.Equal(t, "x", optionalString("x").ValueString())
}

// Every registered definition needs an engine and a settings field in the resource model, or the resource cannot
// handle it.
func TestIntegrationDefinitionsAreWired(t *testing.T) {
	for _, d := range integrationDefinitions {
		t.Run(d.Name(), func(t *testing.T) {
			assert.True(t, isMultiInstance(d), "no engine handles %s", d.Name())
			m := integrationConfigurationModel{}
			assert.NotNil(t, m.settings(d.Name()), "integrationConfigurationModel has no settings field for %s", d.Name())
		})
	}
}

// Every credential kind needs a model field and a part() case for each of its parts, or drift detection silently
// skips the missing part.
func TestCredentialKindsAreWired(t *testing.T) {
	models := map[credentialKind]*credentialsModel{
		credentialToken:   {Token: &tokenCredentialModel{}},
		credentialBasic:   {Basic: &basicCredentialModel{}},
		credentialKeyPair: {KeyPair: &keyPairCredentialModel{}},
	}
	assert.Len(t, models, len(credentialParts), "every kind in credentialParts needs a model here")
	for kind, parts := range credentialParts {
		c := models[kind]
		require.NotNil(t, c, "no model for kind %s", kind)
		assert.Equal(t, kind, c.kind())
		for _, p := range parts {
			assert.NotNil(t, c.part(p.name), "kind %s has no part() case for %s", kind, p.name)
		}
		schemaKind, ok := credentialsAttribute().Attributes[string(kind)].(schema.SingleNestedAttribute)
		require.True(t, ok, "credentials schema has no block for %s", kind)
		assert.Len(t, schemaKind.Attributes, len(parts), "credentials.%s schema and credentialParts differ", kind)
		for _, p := range parts {
			attribute, ok := schemaKind.Attributes[p.name].(schema.StringAttribute)
			require.True(t, ok, "credentials.%s schema has no %s", kind, p.name)
			assert.Equal(t, p.secret, attribute.Sensitive, "credentials.%s.%s sensitivity differs from credentialParts", kind, p.name)
		}
	}
}
