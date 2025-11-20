package provider

import (
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestCatalogEntityOwnerResourceModel_ToApiModel_WithInheritance(t *testing.T) {
	tests := []struct {
		name     string
		model    CatalogEntityOwnerResourceModel
		expected cortex.CatalogEntityOwner
	}{
		{
			name: "Email owner with APPEND inheritance",
			model: CatalogEntityOwnerResourceModel{
				Type:        types.StringValue("EMAIL"),
				Name:        types.StringValue("John Doe"),
				Email:       types.StringValue("john.doe@example.com"),
				Description: types.StringValue("Engineering Lead"),
				Inheritance: types.StringValue("APPEND"),
			},
			expected: cortex.CatalogEntityOwner{
				Type:        "EMAIL",
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Description: "Engineering Lead",
				Inheritance: "APPEND",
			},
		},
		{
			name: "Group owner with FALLBACK inheritance",
			model: CatalogEntityOwnerResourceModel{
				Type:        types.StringValue("GROUP"),
				Name:        types.StringValue("SRE Team"),
				Provider:    types.StringValue("CORTEX"),
				Description: types.StringValue("Site Reliability Engineering"),
				Inheritance: types.StringValue("FALLBACK"),
			},
			expected: cortex.CatalogEntityOwner{
				Type:        "GROUP",
				Name:        "SRE Team",
				Provider:    "CORTEX",
				Description: "Site Reliability Engineering",
				Inheritance: "FALLBACK",
			},
		},
		{
			name: "Slack owner with NONE inheritance",
			model: CatalogEntityOwnerResourceModel{
				Type:                 types.StringValue("SLACK"),
				Channel:              types.StringValue("engineering"),
				Description:          types.StringValue("Engineering channel"),
				NotificationsEnabled: types.BoolValue(true),
				Inheritance:          types.StringValue("NONE"),
			},
			expected: cortex.CatalogEntityOwner{
				Type:                 "SLACK",
				Channel:              "engineering",
				Description:          "Engineering channel",
				NotificationsEnabled: true,
				Inheritance:          "NONE",
			},
		},
		{
			name: "Owner without inheritance (null)",
			model: CatalogEntityOwnerResourceModel{
				Type:        types.StringValue("EMAIL"),
				Name:        types.StringValue("Jane Doe"),
				Email:       types.StringValue("jane.doe@example.com"),
				Inheritance: types.StringNull(),
			},
			expected: cortex.CatalogEntityOwner{
				Type:        "EMAIL",
				Name:        "Jane Doe",
				Email:       "jane.doe@example.com",
				Inheritance: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.model.ToApiModel()
			assert.Equal(t, tt.expected.Type, actual.Type)
			assert.Equal(t, tt.expected.Name, actual.Name)
			assert.Equal(t, tt.expected.Email, actual.Email)
			assert.Equal(t, tt.expected.Description, actual.Description)
			assert.Equal(t, tt.expected.Provider, actual.Provider)
			assert.Equal(t, tt.expected.Channel, actual.Channel)
			assert.Equal(t, tt.expected.NotificationsEnabled, actual.NotificationsEnabled)
			assert.Equal(t, tt.expected.Inheritance, actual.Inheritance)
		})
	}
}

func TestCatalogEntityOwnerResourceModel_FromApiModel_WithInheritance(t *testing.T) {
	tests := []struct {
		name     string
		apiModel cortex.CatalogEntityOwner
		expected CatalogEntityOwnerResourceModel
	}{
		{
			name: "Email owner with APPEND inheritance",
			apiModel: cortex.CatalogEntityOwner{
				Type:        "email",
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Description: "Engineering Lead",
				Inheritance: "APPEND",
			},
			expected: CatalogEntityOwnerResourceModel{
				Type:                 types.StringValue("email"),
				Name:                 types.StringValue("John Doe"),
				Email:                types.StringValue("john.doe@example.com"),
				Description:          types.StringValue("Engineering Lead"),
				Inheritance:          types.StringValue("APPEND"),
				Channel:              types.StringNull(),
				Provider:             types.StringNull(),
				NotificationsEnabled: types.BoolNull(),
			},
		},
		{
			name: "Group owner with FALLBACK inheritance",
			apiModel: cortex.CatalogEntityOwner{
				Type:        "group",
				Name:        "SRE Team",
				Provider:    "CORTEX",
				Description: "Site Reliability Engineering",
				Inheritance: "FALLBACK",
			},
			expected: CatalogEntityOwnerResourceModel{
				Type:                 types.StringValue("group"),
				Name:                 types.StringValue("SRE Team"),
				Provider:             types.StringValue("CORTEX"),
				Description:          types.StringValue("Site Reliability Engineering"),
				Inheritance:          types.StringValue("FALLBACK"),
				Email:                types.StringNull(),
				Channel:              types.StringNull(),
				NotificationsEnabled: types.BoolNull(),
			},
		},
		{
			name: "Slack owner with NONE inheritance",
			apiModel: cortex.CatalogEntityOwner{
				Type:                 "slack",
				Channel:              "engineering",
				Description:          "Engineering channel",
				NotificationsEnabled: true,
				Inheritance:          "NONE",
			},
			expected: CatalogEntityOwnerResourceModel{
				Type:                 types.StringValue("slack"),
				Channel:              types.StringValue("engineering"),
				Description:          types.StringValue("Engineering channel"),
				NotificationsEnabled: types.BoolValue(true),
				Inheritance:          types.StringValue("NONE"),
				Email:                types.StringNull(),
				Name:                 types.StringNull(),
				Provider:             types.StringNull(),
			},
		},
		{
			name: "Owner without inheritance (empty string)",
			apiModel: cortex.CatalogEntityOwner{
				Type:        "email",
				Name:        "Jane Doe",
				Email:       "jane.doe@example.com",
				Inheritance: "",
			},
			expected: CatalogEntityOwnerResourceModel{
				Type:                 types.StringValue("email"),
				Name:                 types.StringValue("Jane Doe"),
				Email:                types.StringValue("jane.doe@example.com"),
				Description:          types.StringNull(),
				Inheritance:          types.StringNull(),
				Channel:              types.StringNull(),
				Provider:             types.StringNull(),
				NotificationsEnabled: types.BoolNull(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &CatalogEntityOwnerResourceModel{}
			actual := model.FromApiModel(&tt.apiModel)

			assert.Equal(t, tt.expected.Type, actual.Type)
			assert.Equal(t, tt.expected.Name, actual.Name)
			assert.Equal(t, tt.expected.Email, actual.Email)
			assert.Equal(t, tt.expected.Description, actual.Description)
			assert.Equal(t, tt.expected.Provider, actual.Provider)
			assert.Equal(t, tt.expected.Channel, actual.Channel)
			assert.Equal(t, tt.expected.NotificationsEnabled, actual.NotificationsEnabled)
			assert.Equal(t, tt.expected.Inheritance, actual.Inheritance)
		})
	}
}

func TestCatalogEntityOwnerResourceModel_RoundTrip_WithInheritance(t *testing.T) {
	tests := []struct {
		name  string
		model CatalogEntityOwnerResourceModel
	}{
		{
			name: "Email with APPEND",
			model: CatalogEntityOwnerResourceModel{
				Type:        types.StringValue("email"),
				Name:        types.StringValue("John Doe"),
				Email:       types.StringValue("john@example.com"),
				Inheritance: types.StringValue("APPEND"),
			},
		},
		{
			name: "Group with FALLBACK",
			model: CatalogEntityOwnerResourceModel{
				Type:        types.StringValue("group"),
				Name:        types.StringValue("Team"),
				Provider:    types.StringValue("CORTEX"),
				Inheritance: types.StringValue("FALLBACK"),
			},
		},
		{
			name: "Slack with NONE",
			model: CatalogEntityOwnerResourceModel{
				Type:                 types.StringValue("slack"),
				Channel:              types.StringValue("general"),
				NotificationsEnabled: types.BoolValue(false),
				Inheritance:          types.StringValue("NONE"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to API model
			apiModel := tt.model.ToApiModel()

			// Convert back to resource model
			resourceModel := &CatalogEntityOwnerResourceModel{}
			actual := resourceModel.FromApiModel(&apiModel)

			// Verify inheritance field is preserved
			assert.Equal(t, tt.model.Inheritance, actual.Inheritance, "Inheritance field should be preserved in round trip")
		})
	}
}
