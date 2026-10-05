package provider

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

var mendTypes = map[string]attr.Type{
	"application_ids": types.ListType{ElemType: types.StringType},
	"project_ids":     types.ListType{ElemType: types.StringType},
}

func mendObject(applicationIds, projectIds types.List) types.Object {
	return types.ObjectValueMust(mendTypes, map[string]attr.Value{"application_ids": applicationIds, "project_ids": projectIds})
}

func stringList(values ...string) types.List {
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elems)
}

func TestKeepEmptyObject(t *testing.T) {
	ctx := context.Background()
	nullList := types.ListNull(types.StringType)
	tests := []struct {
		name  string
		prior types.Object
		read  types.Object
		want  types.Object
	}{
		{"empty prior, null read", mendObject(stringList(), stringList()), types.ObjectNull(mendTypes), mendObject(stringList(), stringList())},
		{"null prior, null read", types.ObjectNull(mendTypes), types.ObjectNull(mendTypes), types.ObjectNull(mendTypes)},
		{"null prior, set read", types.ObjectNull(mendTypes), mendObject(stringList("1"), nullList), mendObject(stringList("1"), nullList)},
		{"empty prior, set read", mendObject(stringList(), stringList()), mendObject(stringList("1"), nullList), mendObject(stringList("1"), stringList())},
		{"set prior, null read", mendObject(stringList("1"), nullList), types.ObjectNull(mendTypes), types.ObjectNull(mendTypes)},
		{"set prior, other read", mendObject(stringList("1"), nullList), mendObject(stringList("2"), nullList), mendObject(stringList("2"), nullList)},
		{"unknown prior", types.ObjectUnknown(mendTypes), types.ObjectNull(mendTypes), types.ObjectNull(mendTypes)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, keepEmptyObject(ctx, tt.prior, tt.read))
		})
	}
}

func TestCatalogEntityStaticAnalysisMendResourceModel_FromApiModel_MissingIdsAreNull(t *testing.T) {
	ctx := context.Background()
	diags := diag.Diagnostics{}
	model := CatalogEntityStaticAnalysisMendResourceModel{}
	got := model.FromApiModel(ctx, &diags, &cortex.CatalogEntityStaticAnalysisMend{ApplicationIDs: []string{"1"}})
	assert.False(t, diags.HasError())
	assert.Equal(t, mendObject(stringList("1"), types.ListNull(types.StringType)), got)
}
