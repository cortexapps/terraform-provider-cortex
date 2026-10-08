# Entity Type icon_tag + Relationship Types — Design Spec

**Date:** 2026-10-08
**Linear:** PFR-4258 — Kraken Tech Terraform support for entity and relationship types

---

## Summary

Two changes:

1. **Add `icon_tag` to `cortex_resource_definition`** — the existing resource already covers entity types; it just needs this one missing optional field.
2. **Add `cortex_relationship_type`** — a new resource and data source for managing relationship types via the Cortex API.

A third component adds MLB-themed test data to the integration test project (`https://github.com/jeff-test-org/cortex-terraform-mlb`) to exercise both changes end-to-end.

---

## Part 1: `icon_tag` on `cortex_resource_definition`

### What changes

The Cortex entity-type API (`/api/v1/catalog/definitions`) already supports an `iconTag` field; it is just not exposed in the provider today.

**`internal/cortex/resource_definitions.go`**
- Add `IconTag string` to `ResourceDefinition`, `CreateResourceDefinitionRequest`, and `UpdateResourceDefinitionRequest` structs (all `json:"iconTag,omitempty"`).
- Add to `ToCreateRequest()` and `ToUpdateRequest()` conversion helpers.

**`internal/provider/resource_definition_models.go`**
- Add `IconTag types.String` with `tfsdk:"icon_tag"` to `ResourceDefinitionResourceModel`.
- Update `FromApiModel()` and `ToApiModel()` to include `icon_tag`.

**`internal/provider/resource_definition_resource.go`**
- Add `icon_tag` to the schema as `schema.StringAttribute{Optional: true}`.

**`internal/provider/resource_definition_data_source.go` + `_models.go`**
- Mirror the same `icon_tag` addition to the data source schema and model.

**`examples/resources/cortex_resource_definition/`** — update example to show `icon_tag`.

### Constraints
- Fully backwards-compatible: existing configs with no `icon_tag` continue to work unchanged.
- No plan modifiers needed — the field is purely optional and the API returns it only when set.

---

## Part 2: `cortex_relationship_type` resource + data source

### API

Base URL: `https://api.getcortexapp.com/api/v1/relationship-types`

| Operation | Method | Path |
|-----------|--------|------|
| Create | POST | `/api/v1/relationship-types` |
| Get | GET | `/api/v1/relationship-types/{tag}` |
| List | GET | `/api/v1/relationship-types` |
| Update | PUT | `/api/v1/relationship-types/{tag}` |
| Delete | DELETE | `/api/v1/relationship-types/{tag}` |

### API model (`EntityRelationshipType`)

```go
type EntityRelationshipType struct {
    Tag                      string              `json:"tag"`
    Name                     string              `json:"name"`
    Description              string              `json:"description,omitempty"`
    AllowCycles              bool                `json:"allowCycles"`
    IsCortexManaged          bool                `json:"isCortexManaged"`
    IsSingleSource           bool                `json:"isSingleSource"`
    IsSingleDestination      bool                `json:"isSingleDestination"`
    DefinitionLocation       string              `json:"definitionLocation"` // SOURCE | DESTINATION | BOTH
    SourcesFilter            EntityTypeFilter    `json:"sourcesFilter"`
    DestinationsFilter       EntityTypeFilter    `json:"destinationsFilter"`
    SourceLabelSingular      string              `json:"sourceLabelSingular,omitempty"`
    SourceLabelPlural        string              `json:"sourceLabelPlural,omitempty"`
    DestinationLabelSingular string              `json:"destinationLabelSingular,omitempty"`
    DestinationLabelPlural   string              `json:"destinationLabelPlural,omitempty"`
    Inheritances             []RelationshipInheritance `json:"inheritances"`
}

type EntityTypeFilter struct {
    Include   bool     `json:"include"`
    Types     []string `json:"types,omitempty"`
    Providers []string `json:"providers,omitempty"` // AWS | GCP | AZURE
}

type RelationshipInheritance struct {
    ParameterTag    string `json:"parameterTag"`    // currently only "x-cortex-owners"
    InheritanceType string `json:"inheritanceType"` // NONE | FALLBACK | APPEND | CONFIGURE_AT_SOURCE
}
```

### Terraform schema (`cortex_relationship_type`)

| Attribute | Type | Required/Optional/Computed | Notes |
|-----------|------|---------------------------|-------|
| `tag` | String | Required | Import ID; API identifier |
| `name` | String | Required | |
| `description` | String | Optional | |
| `allow_cycles` | Bool | Required | |
| `is_single_source` | Bool | Required | |
| `is_single_destination` | Bool | Required | |
| `definition_location` | String | Required | SOURCE, DESTINATION, or BOTH |
| `sources_filter` | SingleNestedAttribute | Required | see below |
| `destinations_filter` | SingleNestedAttribute | Required | see below |
| `source_label_singular` | String | Optional | |
| `source_label_plural` | String | Optional | |
| `destination_label_singular` | String | Optional | |
| `destination_label_plural` | String | Optional | |
| `inheritances` | ListNestedAttribute | Optional | see below |
| `is_cortex_managed` | Bool | Computed | UseStateForUnknown |

**`sources_filter` / `destinations_filter` nested attributes:**

| Attribute | Type | Required/Optional |
|-----------|------|-------------------|
| `include` | Bool | Required |
| `types` | SetAttribute(String) | Optional |
| `providers` | SetAttribute(String) | Optional |

`types` and `providers` use `SetAttribute` (not List) because order is not guaranteed by the API.

**`inheritances` list item:**

| Attribute | Type | Required/Optional |
|-----------|------|-------------------|
| `parameter_tag` | String | Required |
| `inheritance_type` | String | Required |

### ImportState behavior

If `is_cortex_managed` is `true`, reject the import with a clear error message explaining the resource is platform-managed and cannot be updated or deleted via the API. This prevents users from getting trapped with unmanageable state.

### New files

**API client:**
- `internal/cortex/relationship_types.go` — `RelationshipTypesClientInterface`, `RelationshipTypesClient`, all CRUD methods
- Register `RelationshipTypes()` method on `HttpClient` in `http_client.go`
- Add `"relationship_types": "/api/v1/relationship-types/"` to `BaseUris`

**Provider:**
- `internal/provider/relationship_type_resource.go`
- `internal/provider/relationship_type_models.go`
- `internal/provider/relationship_type_data_source.go`
- `internal/provider/relationship_type_data_source_models.go`
- Register in `provider.go` `Resources()` and `DataSources()`

**Examples:**
- `examples/resources/cortex_relationship_type/resource.tf`
- `examples/data-sources/cortex_relationship_type/data-source.tf`

**Acceptance tests:**
- `internal/provider/relationship_type_resource_test.go` — minimal create+read+import, full config with optional fields

---

## Part 3: MLB integration test additions

Repository: `https://github.com/jeff-test-org/cortex-terraform-mlb`

### Entity type

```hcl
resource "cortex_resource_definition" "mlb_leader" {
  type        = "mlb-leader"
  name        = "MLB Statistical Leader"
  description = "Top 10 leaders for an MLB statistical category in a given season"
  schema = jsonencode({
    properties = {
      category = { type = "string", description = "Statistical category (HR, RBI, AVG, etc.)" }
      season   = { type = "integer", description = "MLB season year" }
      league   = { type = "string", enum = ["AL", "NL"] }
      leaders  = {
        type  = "array"
        items = {
          type       = "object"
          properties = {
            rank   = { type = "integer" }
            player = { type = "string" }
            team   = { type = "string" }
            value  = { type = "string" }
          }
        }
      }
    }
  })
}
```

### Relationship type

```hcl
resource "cortex_relationship_type" "mlb_to_leader" {
  tag  = "mlb-to-leader"
  name = "MLB to Statistical Leader"
  definition_location  = "SOURCE"
  allow_cycles         = false
  is_single_source     = false
  is_single_destination = false
  sources_filter = {
    include = true
    types   = ["mlb"]
  }
  destinations_filter = {
    include = true
    types   = ["mlb-leader"]
  }
  inheritances = []
}
```

### Leader entities (12 total)

Six AL categories: `2026-al-leader-hr`, `2026-al-leader-rbi`, `2026-al-leader-avg`, `2026-al-leader-wins`, `2026-al-leader-era`, `2026-al-leader-strikeouts`

Six NL categories: same pattern with `nl`.

All created as `cortex_catalog_entity` resources of type `mlb-leader`.

### Relationships

Relationships from the existing `2026-al` entity to each of the 6 AL leader entities, and from `2026-nl` to each of the 6 NL leader entities, using the `mlb-to-leader` relationship type. Implementation TBD based on what `cortex_catalog_entity` supports for relationships (may require a separate relationship resource if the provider exposes one).

---

## Success criteria

- `terraform plan` on existing configs that use `cortex_resource_definition` shows no diff (no breaking change)
- `icon_tag` round-trips correctly: set → apply → plan shows no diff
- `cortex_relationship_type` create, read, update, delete all work against the real API
- Importing a Cortex-managed relationship type produces a clear error
- MLB test project applies cleanly and entities/relationship types are visible in the Cortex UI
- `go generate ./...` produces clean docs
- `golangci-lint` passes
