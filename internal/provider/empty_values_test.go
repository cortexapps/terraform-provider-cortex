package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

var mendTypes = (&CatalogEntityStaticAnalysisMendResourceModel{}).AttrTypes()

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

func TestKeepEmptyObject_Nested(t *testing.T) {
	ctx := context.Background()
	outerTypes := map[string]attr.Type{"mend": types.ObjectType{AttrTypes: mendTypes}, "name": types.StringType}
	outer := func(mend types.Object, name types.String) types.Object {
		return types.ObjectValueMust(outerTypes, map[string]attr.Value{"mend": mend, "name": name})
	}
	emptyMend := mendObject(stringList(), types.ListNull(types.StringType))

	t.Run("keeps an empty child of a set parent", func(t *testing.T) {
		got := keepEmptyObject(ctx, outer(emptyMend, types.StringValue("a")), outer(types.ObjectNull(mendTypes), types.StringValue("a")))
		assert.Equal(t, outer(emptyMend, types.StringValue("a")), got)
	})
	t.Run("reads a changed sibling", func(t *testing.T) {
		got := keepEmptyObject(ctx, outer(emptyMend, types.StringValue("a")), outer(types.ObjectNull(mendTypes), types.StringValue("b")))
		assert.Equal(t, outer(emptyMend, types.StringValue("b")), got)
	})
	t.Run("reads a child set remotely", func(t *testing.T) {
		setMend := mendObject(stringList("1"), types.ListNull(types.StringType))
		got := keepEmptyObject(ctx, outer(emptyMend, types.StringNull()), outer(setMend, types.StringNull()))
		assert.Equal(t, outer(mendObject(stringList("1"), types.ListNull(types.StringType)), types.StringNull()), got)
	})
}
