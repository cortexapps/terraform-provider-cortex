package provider

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogEntityDataSourceModel_FromApiModel_HappyPath(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "A test service",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{"platform", "backend"},
		Links: []cortex.CatalogEntityLink{
			{
				Name:        "runbook",
				Type:        "runbook",
				Url:         "https://runbook.example.com",
				Description: "Service runbook",
			},
			{
				Name:        "dashboard",
				Type:        "dashboard",
				Url:         "https://dashboard.example.com",
				Description: "Service metrics",
			},
		},
		Metadata: []cortex.CatalogEntityMetadata{
			{
				Key:   "team",
				Value: "platform-team",
			},
			{
				Key:   "cost-center",
				Value: 12345,
			},
		},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{
				{
					Email:       "owner@example.com",
					Description: "Primary owner",
					Inheritance: "DIRECT",
				},
			},
			Groups: []cortex.CatalogEntityGroup{
				{
					GroupName:   "platform-team",
					Description: "Platform team",
					Provider:    "okta",
					Inheritance: "DIRECT",
				},
			},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{
			{
				Name:                 "platform-alerts",
				NotificationsEnabled: true,
				Description:          "Alerts channel",
			},
		},
		Git: cortex.CatalogEntityGitSummary{
			Provider:      "github",
			Repository:    "org/repo",
			Alias:         "main-repo",
			Basepath:      "/services/test-service",
			RepositoryUrl: "https://github.com/org/repo",
		},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify simple fields
	assert.Equal(t, "test-service", model.Id.ValueString())
	assert.Equal(t, "test-service", model.Tag.ValueString())
	assert.Equal(t, "Test Service", model.Name.ValueString())
	assert.Equal(t, "A test service", model.Description.ValueString())
	assert.Equal(t, "service", model.Type.ValueString())
	assert.False(t, model.IsArchived.ValueBool())
	assert.Equal(t, "2024-01-15T10:30:00Z", model.LastUpdated.ValueString())

	// Verify groups
	require.Len(t, model.Groups, 2)
	assert.Equal(t, "platform", model.Groups[0].ValueString())
	assert.Equal(t, "backend", model.Groups[1].ValueString())

	// Verify links
	require.Len(t, model.Links, 2)
	assert.Equal(t, "runbook", model.Links[0].Name.ValueString())
	assert.Equal(t, "runbook", model.Links[0].Type.ValueString())
	assert.Equal(t, "https://runbook.example.com", model.Links[0].Url.ValueString())
	assert.Equal(t, "Service runbook", model.Links[0].Description.ValueString())

	// Verify metadata (complex values are JSON serialized)
	require.Len(t, model.Metadata, 2)
	assert.Equal(t, "team", model.Metadata[0].Key.ValueString())
	assert.Equal(t, `"platform-team"`, model.Metadata[0].Value.ValueString())
	assert.Equal(t, "cost-center", model.Metadata[1].Key.ValueString())
	assert.Equal(t, `12345`, model.Metadata[1].Value.ValueString())

	// Verify ownership collection is populated
	require.False(t, model.Ownership.IsNull())
	var ownershipVal CatalogEntityOwnershipModel
	ownershipDiags := model.Ownership.As(ctx, &ownershipVal, basetypes.ObjectAsOptions{})
	require.False(t, ownershipDiags.HasError())
	require.Len(t, ownershipVal.Emails, 1)
	require.Len(t, ownershipVal.Groups, 1)

	// Verify slack channels
	require.Len(t, model.SlackChannels, 1)
	var sc0 CatalogEntitySlackChannelModel
	sc0Diags := model.SlackChannels[0].As(ctx, &sc0, basetypes.ObjectAsOptions{})
	require.False(t, sc0Diags.HasError())
	assert.Equal(t, "platform-alerts", sc0.Name.ValueString())
	assert.True(t, sc0.NotificationsEnabled.ValueBool())

	// Verify git
	require.False(t, model.Git.IsNull())
	var gitVal CatalogEntityGitModel
	gitDiags := model.Git.As(ctx, &gitVal, basetypes.ObjectAsOptions{})
	require.False(t, gitDiags.HasError())
	assert.Equal(t, "github", gitVal.Provider.ValueString())
	assert.Equal(t, "org/repo", gitVal.Repository.ValueString())
	assert.Equal(t, "main-repo", gitVal.Alias.ValueString())
	assert.Equal(t, "/services/test-service", gitVal.Basepath.ValueString())
	assert.Equal(t, "https://github.com/org/repo", gitVal.RepositoryUrl.ValueString())
}

func TestCatalogEntityDataSourceModel_FromApiModel_EmptyCollections(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "minimal-service",
		Name:        "Minimal Service",
		Description: "",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git: cortex.CatalogEntityGitSummary{
			Provider:      "",
			Repository:    "",
			Alias:         "",
			Basepath:      "",
			RepositoryUrl: "",
		},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify empty collections are properly initialized
	assert.Len(t, model.Groups, 0)
	assert.Len(t, model.Links, 0)
	assert.Len(t, model.Metadata, 0)
	assert.Len(t, model.SlackChannels, 0)

	// Verify ownership object exists even with empty collections
	require.False(t, model.Ownership.IsNull())
	var ownershipVal CatalogEntityOwnershipModel
	ownershipDiags := model.Ownership.As(ctx, &ownershipVal, basetypes.ObjectAsOptions{})
	require.False(t, ownershipDiags.HasError())
	assert.Len(t, ownershipVal.Emails, 0)
	assert.Len(t, ownershipVal.Groups, 0)
}

func TestCatalogEntityDataSourceModel_FromApiModel_MetadataWithComplexValues(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "A test",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata: []cortex.CatalogEntityMetadata{
			{
				Key:   "string-value",
				Value: "string-data",
			},
			{
				Key:   "number-value",
				Value: 42,
			},
			{
				Key:   "float-value",
				Value: 3.14,
			},
			{
				Key:   "boolean-value",
				Value: true,
			},
			{
				Key: "object-value",
				Value: map[string]interface{}{
					"nested": "data",
					"count":  10,
				},
			},
			{
				Key:   "array-value",
				Value: []interface{}{"item1", "item2"},
			},
		},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git:           cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify metadata values are properly JSON serialized
	require.Len(t, model.Metadata, 6)
	assert.Equal(t, `"string-data"`, model.Metadata[0].Value.ValueString())
	assert.Equal(t, `42`, model.Metadata[1].Value.ValueString())
	assert.Equal(t, `3.14`, model.Metadata[2].Value.ValueString())
	assert.Equal(t, `true`, model.Metadata[3].Value.ValueString())
	assert.Contains(t, model.Metadata[4].Value.ValueString(), "nested")
	assert.Contains(t, model.Metadata[5].Value.ValueString(), "item1")
}

func TestCatalogEntityDataSourceModel_FromApiModel_MultipleOwnershipEmails(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{
				{Email: "owner1@example.com", Description: "Primary", Inheritance: "DIRECT"},
				{Email: "owner2@example.com", Description: "Secondary", Inheritance: "INDIRECT"},
				{Email: "owner3@example.com", Description: "Backup", Inheritance: "DIRECT"},
			},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git:           cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify ownership emails are populated with correct values
	require.False(t, model.Ownership.IsNull())
	var ownershipVal CatalogEntityOwnershipModel
	ownershipDiags := model.Ownership.As(ctx, &ownershipVal, basetypes.ObjectAsOptions{})
	require.False(t, ownershipDiags.HasError())
	require.Len(t, ownershipVal.Emails, 3)
	var email0 CatalogEntityOwnershipEmailModel
	email0Diags := ownershipVal.Emails[0].As(ctx, &email0, basetypes.ObjectAsOptions{})
	require.False(t, email0Diags.HasError())
	assert.Equal(t, "owner1@example.com", email0.Email.ValueString())
	assert.Equal(t, "Primary", email0.Description.ValueString())
	assert.Equal(t, "DIRECT", email0.Inheritance.ValueString())
	var email1 CatalogEntityOwnershipEmailModel
	email1Diags := ownershipVal.Emails[1].As(ctx, &email1, basetypes.ObjectAsOptions{})
	require.False(t, email1Diags.HasError())
	assert.Equal(t, "owner2@example.com", email1.Email.ValueString())
	assert.Equal(t, "INDIRECT", email1.Inheritance.ValueString())
}

func TestCatalogEntityDataSourceModel_FromApiModel_MultipleOwnershipGroups(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{
				{GroupName: "team-a", Description: "Team A", Provider: "okta", Inheritance: "DIRECT"},
				{GroupName: "team-b", Description: "Team B", Provider: "azure", Inheritance: "INDIRECT"},
				{GroupName: "team-c", Description: "Team C", Provider: "okta", Inheritance: "DIRECT"},
			},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git:           cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify ownership groups are populated with correct values
	require.False(t, model.Ownership.IsNull())
	var ownershipVal CatalogEntityOwnershipModel
	ownershipDiags := model.Ownership.As(ctx, &ownershipVal, basetypes.ObjectAsOptions{})
	require.False(t, ownershipDiags.HasError())
	require.Len(t, ownershipVal.Groups, 3)
	var group0 CatalogEntityOwnershipGroupModel
	group0Diags := ownershipVal.Groups[0].As(ctx, &group0, basetypes.ObjectAsOptions{})
	require.False(t, group0Diags.HasError())
	assert.Equal(t, "team-a", group0.GroupName.ValueString())
	assert.Equal(t, "Team A", group0.Description.ValueString())
	assert.Equal(t, "okta", group0.Provider.ValueString())
	assert.Equal(t, "DIRECT", group0.Inheritance.ValueString())
	var group1 CatalogEntityOwnershipGroupModel
	group1Diags := ownershipVal.Groups[1].As(ctx, &group1, basetypes.ObjectAsOptions{})
	require.False(t, group1Diags.HasError())
	assert.Equal(t, "team-b", group1.GroupName.ValueString())
	assert.Equal(t, "azure", group1.Provider.ValueString())
	assert.Equal(t, "INDIRECT", group1.Inheritance.ValueString())
}

func TestCatalogEntityDataSourceModel_FromApiModel_MultipleSlackChannels(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{
			{Name: "alerts", NotificationsEnabled: true, Description: "Alert notifications"},
			{Name: "updates", NotificationsEnabled: false, Description: "Update announcements"},
			{Name: "general", NotificationsEnabled: true, Description: "General discussion"},
		},
		Git: cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify slack channel fields are correctly mapped
	require.Len(t, model.SlackChannels, 3)
	var sc0 CatalogEntitySlackChannelModel
	sc0Diags := model.SlackChannels[0].As(ctx, &sc0, basetypes.ObjectAsOptions{})
	require.False(t, sc0Diags.HasError())
	assert.Equal(t, "alerts", sc0.Name.ValueString())
	assert.True(t, sc0.NotificationsEnabled.ValueBool())
	assert.Equal(t, "Alert notifications", sc0.Description.ValueString())
	var sc1 CatalogEntitySlackChannelModel
	sc1Diags := model.SlackChannels[1].As(ctx, &sc1, basetypes.ObjectAsOptions{})
	require.False(t, sc1Diags.HasError())
	assert.Equal(t, "updates", sc1.Name.ValueString())
	assert.False(t, sc1.NotificationsEnabled.ValueBool())
}

func TestCatalogEntityDataSourceModel_FromApiModel_ArchivedEntity(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "archived-service",
		Name:        "Archived Service",
		Description: "An archived service",
		Type:        "service",
		IsArchived:  true,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git:           cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify is_archived is properly set to true
	assert.True(t, model.IsArchived.ValueBool())
}

func TestCatalogEntityDataSourceModel_FromApiModel_VariousEntityTypes(t *testing.T) {
	ctx := context.Background()

	entityTypes := []string{"service", "resource", "domain", "team"}

	for _, entityType := range entityTypes {
		t.Run(entityType, func(t *testing.T) {
			diagnostics := diag.Diagnostics{}

			entity := &cortex.CatalogEntity{
				Tag:         "test-" + entityType,
				Name:        "Test " + entityType,
				Description: "A test " + entityType,
				Type:        entityType,
				IsArchived:  false,
				LastUpdated: "2024-01-15T10:30:00Z",
				Groups:      []string{},
				Links:       []cortex.CatalogEntityLink{},
				Metadata:    []cortex.CatalogEntityMetadata{},
				Ownership: cortex.CatalogEntityOwnership{
					Emails: []cortex.CatalogEntityEmail{},
					Groups: []cortex.CatalogEntityGroup{},
				},
				SlackChannels: []cortex.CatalogEntitySlackChannel{},
				Git:           cortex.CatalogEntityGitSummary{},
			}

			model := &CatalogEntityDataSourceModel{}
			model.FromApiModel(ctx, &diagnostics, entity)

			// Verify no diagnostics errors
			require.False(t, diagnostics.HasError())

			// Verify type is correctly mapped
			assert.Equal(t, entityType, model.Type.ValueString())
		})
	}
}

func TestCatalogEntityDataSourceModel_FromApiModel_GitInformation(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git: cortex.CatalogEntityGitSummary{
			Provider:      "github",
			Repository:    "my-org/my-repo",
			Alias:         "primary",
			Basepath:      "/src/services",
			RepositoryUrl: "https://github.com/my-org/my-repo",
		},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify git fields are correctly mapped
	require.False(t, model.Git.IsNull())
	var gitVal CatalogEntityGitModel
	gitDiags := model.Git.As(ctx, &gitVal, basetypes.ObjectAsOptions{})
	require.False(t, gitDiags.HasError())
	assert.Equal(t, "github", gitVal.Provider.ValueString())
	assert.Equal(t, "my-org/my-repo", gitVal.Repository.ValueString())
	assert.Equal(t, "primary", gitVal.Alias.ValueString())
	assert.Equal(t, "/src/services", gitVal.Basepath.ValueString())
	assert.Equal(t, "https://github.com/my-org/my-repo", gitVal.RepositoryUrl.ValueString())
}

func TestCatalogEntityDataSourceModel_FromApiModel_NilCollections(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	// Entity with nil slices (not empty, but nil)
	entity := &cortex.CatalogEntity{
		Tag:         "minimal",
		Name:        "Minimal",
		Description: "",
		Type:        "",
		IsArchived:  false,
		LastUpdated: "",
		Groups:      nil,
		Links:       nil,
		Metadata:    nil,
		Ownership: cortex.CatalogEntityOwnership{
			Emails: nil,
			Groups: nil,
		},
		SlackChannels: nil,
		Git:           cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Should handle nil collections gracefully without errors
	require.False(t, diagnostics.HasError())
	assert.Equal(t, "minimal", model.Tag.ValueString())
	// Nil slices should be initialized to empty
	assert.Len(t, model.Groups, 0)
	assert.Len(t, model.Links, 0)
	assert.Len(t, model.Metadata, 0)
	assert.Len(t, model.SlackChannels, 0)
}

func TestCatalogEntityDataSourceModel_FromApiModel_MultipleLinks(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "Test Service",
		Description: "",
		Type:        "service",
		IsArchived:  false,
		LastUpdated: "2024-01-15T10:30:00Z",
		Groups:      []string{},
		Links: []cortex.CatalogEntityLink{
			{Name: "runbook", Type: "runbook", Url: "https://example.com/runbook", Description: "Runbook link"},
			{Name: "docs", Type: "documentation", Url: "https://example.com/docs", Description: "Documentation"},
			{Name: "dashboard", Type: "dashboard", Url: "https://example.com/dash", Description: "Metrics dashboard"},
			{Name: "repo", Type: "repository", Url: "https://github.com/org/repo", Description: "Source code"},
		},
		Metadata: []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git:           cortex.CatalogEntityGitSummary{},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify all links are properly converted
	require.Len(t, model.Links, 4)
	assert.Equal(t, "runbook", model.Links[0].Name.ValueString())
	assert.Equal(t, "docs", model.Links[1].Name.ValueString())
	assert.Equal(t, "dashboard", model.Links[2].Name.ValueString())
	assert.Equal(t, "repo", model.Links[3].Name.ValueString())
}

func TestCatalogEntityDataSourceModel_FromApiModel_EmptyStringValues(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}

	entity := &cortex.CatalogEntity{
		Tag:         "test-service",
		Name:        "",
		Description: "",
		Type:        "",
		IsArchived:  false,
		LastUpdated: "",
		Groups:      []string{},
		Links:       []cortex.CatalogEntityLink{},
		Metadata:    []cortex.CatalogEntityMetadata{},
		Ownership: cortex.CatalogEntityOwnership{
			Emails: []cortex.CatalogEntityEmail{},
			Groups: []cortex.CatalogEntityGroup{},
		},
		SlackChannels: []cortex.CatalogEntitySlackChannel{},
		Git: cortex.CatalogEntityGitSummary{
			Provider:      "",
			Repository:    "",
			Alias:         "",
			Basepath:      "",
			RepositoryUrl: "",
		},
	}

	model := &CatalogEntityDataSourceModel{}
	model.FromApiModel(ctx, &diagnostics, entity)

	// Verify no diagnostics errors
	require.False(t, diagnostics.HasError())

	// Verify empty string values are still populated (not null)
	assert.Equal(t, "test-service", model.Tag.ValueString())
	assert.Equal(t, "", model.Name.ValueString())
	assert.Equal(t, "", model.Description.ValueString())
	assert.Equal(t, "", model.Type.ValueString())
	assert.Equal(t, "", model.LastUpdated.ValueString())
}
