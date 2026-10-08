package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

type EntityRelationshipResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	SourceEntityTag       types.String `tfsdk:"source_entity_tag"`
	RelationshipTypeTag   types.String `tfsdk:"relationship_type_tag"`
	DestinationEntityTags types.Set    `tfsdk:"destination_entity_tags"`
}

// FromApiModel populates the model from the API response values.
// Empty destination slice results in an empty (non-null) set, because an empty
// set is a valid and meaningful state for this resource (it means "no destinations").
func (m *EntityRelationshipResourceModel) FromApiModel(sourceTag string, relTypeTag string, destinationTags []string) {
	m.ID = types.StringValue(fmt.Sprintf("%s:%s", sourceTag, relTypeTag))
	m.SourceEntityTag = types.StringValue(sourceTag)
	m.RelationshipTypeTag = types.StringValue(relTypeTag)

	vals := make([]attr.Value, len(destinationTags))
	for i, tag := range destinationTags {
		vals[i] = types.StringValue(tag)
	}
	m.DestinationEntityTags = types.SetValueMust(types.StringType, vals)
}

// ToDestinationTags extracts the destination entity tags as a string slice.
func (m *EntityRelationshipResourceModel) ToDestinationTags() []string {
	return setToStrings(m.DestinationEntityTags)
}
