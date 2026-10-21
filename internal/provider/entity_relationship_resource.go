package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EntityRelationshipResource{}
var _ resource.ResourceWithImportState = &EntityRelationshipResource{}

func NewEntityRelationshipResource() resource.Resource {
	return &EntityRelationshipResource{}
}

func NewEntityRelationshipResourceModel() EntityRelationshipResourceModel {
	return EntityRelationshipResourceModel{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// EntityRelationshipResource defines the resource implementation.
type EntityRelationshipResource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *EntityRelationshipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity_relationship"
}

func (r *EntityRelationshipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the full set of destination entities for a given source entity and relationship type. " +
			"Creating or updating this resource replaces all destinations for the specified source + relationship type combination. " +
			"Deleting this resource clears all destinations (sets them to empty).",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite identifier in the format `source_entity_tag:relationship_type_tag`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_entity_tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the source catalog entity.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"relationship_type_tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the relationship type.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_entity_tags": schema.SetAttribute{
				MarkdownDescription: "Set of entity tags that are destinations for this relationship from the source entity. " +
					"An empty set is valid and clears all destinations.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *EntityRelationshipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cortex.HttpClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *EntityRelationshipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewEntityRelationshipResourceModel()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceTag := data.SourceEntityTag.ValueString()
	relTypeTag := data.RelationshipTypeTag.ValueString()
	destinationTags := data.ToDestinationTags()

	_, err := r.client.EntityRelationships().SetDestinations(ctx, sourceTag, relTypeTag, destinationTags)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create entity relationship, got error: %s", err))
		return
	}

	// The PUT response body may lag (empty destinations) even on a 200, so wait
	// for the GET endpoint instead.
	_, consistent, err := r.getDestinationsConsistent(ctx, sourceTag, relTypeTag, destinationTags)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read entity relationship after create, got error: %s", err))
		return
	}
	if !consistent {
		addStaleDestinationsWarning(&resp.Diagnostics, sourceTag, relTypeTag)
	}

	// The write succeeded, so state holds the applied tags. A read that still
	// lags would make the result disagree with the plan and fail the apply.
	data.FromApiModel(sourceTag, relTypeTag, destinationTags)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EntityRelationshipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewEntityRelationshipResourceModel()

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceTag := data.SourceEntityTag.ValueString()
	relTypeTag := data.RelationshipTypeTag.ValueString()

	// Use the state's destination tags as a hint for the consistency poll.
	// After a write, the GET endpoint may briefly return stale results.
	// We poll until the API returns a result that contains every tag in state.
	stateTags := data.ToDestinationTags()

	destinationTags, _, err := r.getDestinationsConsistent(ctx, sourceTag, relTypeTag, stateTags)
	if err != nil {
		if errors.Is(err, cortex.ApiErrorNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read entity relationship, got error: %s", err))
		return
	}

	data.FromApiModel(sourceTag, relTypeTag, destinationTags)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EntityRelationshipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewEntityRelationshipResourceModel()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceTag := data.SourceEntityTag.ValueString()
	relTypeTag := data.RelationshipTypeTag.ValueString()
	destinationTags := data.ToDestinationTags()

	_, err := r.client.EntityRelationships().SetDestinations(ctx, sourceTag, relTypeTag, destinationTags)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update entity relationship, got error: %s", err))
		return
	}

	_, consistent, err := r.getDestinationsConsistent(ctx, sourceTag, relTypeTag, destinationTags)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read entity relationship after update, got error: %s", err))
		return
	}
	if !consistent {
		addStaleDestinationsWarning(&resp.Diagnostics, sourceTag, relTypeTag)
	}

	data.FromApiModel(sourceTag, relTypeTag, destinationTags)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EntityRelationshipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewEntityRelationshipResourceModel()

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceTag := data.SourceEntityTag.ValueString()
	relTypeTag := data.RelationshipTypeTag.ValueString()

	_, err := r.client.EntityRelationships().SetDestinations(ctx, sourceTag, relTypeTag, []string{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete entity relationship (clear destinations), got error: %s", err))
		return
	}
}

func (r *EntityRelationshipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'source_entity_tag:relationship_type_tag', got: %q", req.ID),
		)
		return
	}

	sourceTag := parts[0]
	relTypeTag := parts[1]

	destinationTags, err := r.client.EntityRelationships().GetDestinations(ctx, sourceTag, relTypeTag)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read entity relationship for import, got error: %s", err))
		return
	}

	data := NewEntityRelationshipResourceModel()
	data.FromApiModel(sourceTag, relTypeTag, destinationTags)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Also set the id attribute explicitly so ImportStatePassthroughID is satisfied.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.ID)...)
}

// getDestinationsConsistent polls GetDestinations to handle the Cortex API's
// brief eventual-consistency window after a write:
//   - When expectedTags is non-empty: poll until the result contains exactly the
//     expected tags (same count, every tag present). Checking both count and
//     membership handles both expansion (stale GET returns too few tags) and
//     shrink (stale GET returns too many — a subset check alone would accept
//     {a,b,c} when the expected set is {a,b}).
//   - When expectedTags is empty: the GET may briefly return the old (non-empty)
//     destination list after a clear. Poll until the result is empty.
//
// Polling terminates as soon as the API reaches the expected state; it is a
// context-aware, bounded poll — not a fixed sleep. The returned bool reports
// whether the API reached the expected state before the wait ran out.
func (r *EntityRelationshipResource) getDestinationsConsistent(ctx context.Context, sourceTag, relTypeTag string, expectedTags []string) ([]string, bool, error) {
	const pollInterval = 250 * time.Millisecond
	maxWait := entityRelationshipConsistencyWait

	consistent := func(tags []string) bool {
		if len(expectedTags) == 0 {
			return len(tags) == 0
		}
		if len(tags) != len(expectedTags) {
			return false
		}
		tagSet := make(map[string]bool, len(tags))
		for _, t := range tags {
			tagSet[t] = true
		}
		for _, e := range expectedTags {
			if !tagSet[e] {
				return false
			}
		}
		return true
	}

	tags, err := r.client.EntityRelationships().GetDestinations(ctx, sourceTag, relTypeTag)
	if err != nil {
		return tags, false, err
	}
	if consistent(tags) {
		return tags, true, nil
	}

	deadline := time.Now().Add(maxWait)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return tags, false, ctx.Err()
		case <-time.After(pollInterval):
		}
		tags, err = r.client.EntityRelationships().GetDestinations(ctx, sourceTag, relTypeTag)
		if err != nil {
			return tags, false, err
		}
		if consistent(tags) {
			return tags, true, nil
		}
	}
	return tags, false, nil
}

// entityRelationshipConsistencyWait bounds how long a write waits for the GET
// endpoint to show it.
var entityRelationshipConsistencyWait = 15 * time.Second

func addStaleDestinationsWarning(diags *diag.Diagnostics, sourceTag, relTypeTag string) {
	diags.AddWarning(
		"Entity relationship not yet visible",
		fmt.Sprintf("Cortex accepted the destinations of %s:%s, but did not return them within %s. "+
			"State holds the applied destinations; the next refresh shows what Cortex stores.",
			sourceTag, relTypeTag, entityRelationshipConsistencyWait),
	)
}
