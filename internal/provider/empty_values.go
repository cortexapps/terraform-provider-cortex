package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// keepEmptyObject returns prior where prior and read are both empty, at any depth. The provider omits an empty block
// or list on write, so without this a configured empty value reads back as null and the apply is inconsistent.
func keepEmptyObject(ctx context.Context, prior, read types.Object) types.Object {
	if isEmptyValue(prior) && isEmptyValue(read) {
		return prior
	}
	if prior.IsNull() || prior.IsUnknown() || read.IsNull() || read.IsUnknown() {
		return read
	}
	priorAttrs := prior.Attributes()
	attrs := map[string]attr.Value{}
	for name, value := range read.Attributes() {
		attrs[name] = value
		switch priorValue := priorAttrs[name].(type) {
		case types.Object:
			if readObject, ok := value.(types.Object); ok {
				attrs[name] = keepEmptyObject(ctx, priorValue, readObject)
			}
		case attr.Value:
			if isEmptyValue(priorValue) && isEmptyValue(value) {
				attrs[name] = priorValue
			}
		}
	}
	// prior and read share the schema's attribute types, so this cannot fail; read is a safe fallback.
	obj, diags := types.ObjectValue(read.AttributeTypes(ctx), attrs)
	if diags.HasError() {
		return read
	}
	return obj
}

// isEmptyValue reports whether v is null, an empty list or set, or an object whose attributes are all empty.
func isEmptyValue(v attr.Value) bool {
	if v.IsNull() {
		return true
	}
	if v.IsUnknown() {
		return false
	}
	switch value := v.(type) {
	case types.Object:
		for _, a := range value.Attributes() {
			if !isEmptyValue(a) {
				return false
			}
		}
		return true
	case types.List:
		return len(value.Elements()) == 0
	case types.Set:
		return len(value.Elements()) == 0
	default:
		return false
	}
}
