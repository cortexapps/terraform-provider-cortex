# Entity Type icon_tag + Relationship Types Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `icon_tag` to `cortex_resource_definition`, add a new `cortex_relationship_type` resource and data source, then validate end-to-end with MLB test data.

**Architecture:** Minimal change to existing resource_definition structs/schema; new relationship_types.go API client and four provider files mirroring the resource_definition pattern. MLB test project updated to use a local provider build for validation.

**Tech Stack:** Go, terraform-plugin-framework v6, terraform-plugin-testing, sling HTTP client.

**Spec:** `docs/superpowers/specs/2026-10-08-entity-and-relationship-types-design.md`

## Global Constraints

- All unordered API collections → `schema.SetAttribute` + `types.Set` (not List)
- Computed-only attributes → `UseStateForUnknown` plan modifier
- `SingleNestedAttribute` examples use `=` attribute syntax, not block syntax
- Run `go generate ./...` after any example file change
- Run `gofmt -w ./...` and `golangci-lint run ./...` before every commit
- `cortex_relationship_type` with `is_cortex_managed = true` → rejected at `ImportState` with clear error
- Acceptance tests follow pattern in `resource_definition_resource_test.go`

## Review Focus

1. **`sources_filter.types` / `providers` returned in different order** — `FromApiModel` must use `types.SetValueMust`, not slice-to-list; test: acceptance test verifies plan is clean after apply.
2. **`is_cortex_managed` shows `(known after apply)` on every update** — covered by `UseStateForUnknown`; test: update step in acceptance test checks no diff.
3. **Importing a Cortex-managed relationship type** — `ImportState` reads the entity and returns an error if `isCortexManaged = true`; not automatically tested in acc tests (requires a real Cortex-managed type), but the guard is in Task 3.
4. **`icon_tag` set then cleared** — when user removes `icon_tag` from config, `ToApiModel` must send `""` not omit the field (otherwise API keeps the old value); checked by an update step that removes `icon_tag`.
5. **`inheritances = []` vs omitted** — empty list and null should both round-trip without diff; covered by minimal (no inheritances) and full (with inheritances) acceptance test cases.

---

## Task 1: Add `icon_tag` to `cortex_resource_definition`

**Files:**
- Modify: `internal/cortex/resource_definitions.go`
- Modify: `internal/provider/resource_definition_models.go`
- Modify: `internal/provider/resource_definition_resource.go`
- Modify: `internal/provider/resource_definition_data_source.go`
- Modify: `internal/provider/resource_definition_data_source_models.go`
- Modify: `examples/resources/resource_definition/resource.tf`
- Test: `internal/provider/resource_definition_resource_test.go`

**Interfaces:**
- Produces: `ResourceDefinition.IconTag string`, `ResourceDefinitionResourceModel.IconTag types.String`

- [ ] **Step 1: Add `IconTag` to API structs**

  In `internal/cortex/resource_definitions.go`:
  - `ResourceDefinition`: add `IconTag string \`json:"iconTag,omitempty"\``
  - `CreateResourceDefinitionRequest`: add `IconTag string \`json:"iconTag,omitempty"\``
  - `UpdateResourceDefinitionRequest`: add `IconTag string \`json:"iconTag,omitempty"\``
  - `ToCreateRequest()`: include `IconTag: r.IconTag`
  - `ToUpdateRequest()`: include `IconTag: r.IconTag`

- [ ] **Step 2: Add `icon_tag` to Terraform resource model**

  In `internal/provider/resource_definition_models.go`:
  - `ResourceDefinitionResourceModel`: add `IconTag types.String \`tfsdk:"icon_tag"\``
  - `FromApiModel`: if `entity.IconTag != ""` → `types.StringValue(entity.IconTag)`, else `types.StringNull()`
  - `ToApiModel`: `IconTag: r.IconTag.ValueString()`

- [ ] **Step 3: Add `icon_tag` to resource schema**

  In `internal/provider/resource_definition_resource.go`, under `// Optional attributes`:
  ```go
  "icon_tag": schema.StringAttribute{
      MarkdownDescription: "Icon tag for the resource definition.",
      Optional:            true,
  },
  ```

- [ ] **Step 4: Add `icon_tag` to data source schema and model**

  In `resource_definition_data_source.go`: add `"icon_tag": schema.StringAttribute{Computed: true}` to schema.
  In `resource_definition_data_source_models.go` (read the file first to find `ResourceDefinitionDataSourceModel`): add `IconTag types.String \`tfsdk:"icon_tag"\`` and update `FromApiModel`.

- [ ] **Step 5: Update example**

  In `examples/resources/resource_definition/resource.tf`, add `icon_tag = "server"` as an optional line in the example resource.

- [ ] **Step 6: Write failing acceptance test for `icon_tag`**

  In `internal/provider/resource_definition_resource_test.go`, add `TestAccResourceDefinitionResourceWithIconTag`:
  - Step 1: create with `icon_tag = "cloud"`; check `resource.TestCheckResourceAttr(..., "icon_tag", "cloud")`
  - Step 2: update removing `icon_tag`; check `resource.TestCheckNoResourceAttr(..., "icon_tag")` (or it equals `""`)
  - Use a unique type like `"test-resourcedef-icontag"`

- [ ] **Step 7: Run acceptance test to verify it fails**

  ```bash
  go test -v ./internal/provider -run TestAccResourceDefinitionResourceWithIconTag
  ```
  Expected: FAIL — schema or API error.

- [ ] **Step 8: Run test to verify it passes**

  (After steps 1–5 are complete.)
  ```bash
  go test -v ./internal/provider -run TestAccResourceDefinitionResourceWithIconTag
  ```
  Expected: PASS.

- [ ] **Step 9: Run `go generate`, lint, format**

  ```bash
  go generate ./...
  gofmt -w ./...
  golangci-lint run ./...
  ```

- [ ] **Step 10: Commit**

  ```bash
  git add internal/cortex/resource_definitions.go \
          internal/provider/resource_definition_models.go \
          internal/provider/resource_definition_resource.go \
          internal/provider/resource_definition_data_source.go \
          internal/provider/resource_definition_data_source_models.go \
          examples/resources/resource_definition/resource.tf \
          internal/provider/resource_definition_resource_test.go \
          docs/
  git commit -m "feat: add icon_tag to cortex_resource_definition"
  ```

---

## Task 2: Relationship Types API client

**Files:**
- Create: `internal/cortex/relationship_types.go`
- Modify: `internal/cortex/http_client.go`

**Interfaces:**
- Produces:
  - `RelationshipTypesClientInterface` with `Get`, `List`, `Create`, `Update`, `Delete`
  - `EntityRelationshipType` struct (all fields from spec)
  - `EntityTypeFilter` struct: `Include bool`, `Types []string`, `Providers []string`
  - `RelationshipInheritance` struct: `ParameterTag string`, `InheritanceType string`
  - `CreateRelationshipTypeRequest` and `UpdateRelationshipTypeRequest`
  - `RelationshipTypesListResponse` with `RelationshipTypes []EntityRelationshipType` and pagination fields
  - `HttpClient.RelationshipTypes() RelationshipTypesClientInterface`

- [ ] **Step 1: Add route to `BaseUris`**

  In `internal/cortex/http_client.go`:
  ```go
  "relationship_types": "/api/v1/relationship-types/",
  ```

- [ ] **Step 2: Create `internal/cortex/relationship_types.go`**

  Define the interface, client struct, and all type structs (see Interfaces above).
  - `Get(ctx, tag string) (EntityRelationshipType, error)` — `GET /api/v1/relationship-types/{tag}`
  - `List(ctx) (RelationshipTypesListResponse, error)` — `GET /api/v1/relationship-types`
  - `Create(ctx, CreateRelationshipTypeRequest) (EntityRelationshipType, error)` — `POST /api/v1/relationship-types`
  - `Update(ctx, tag string, UpdateRelationshipTypeRequest) (EntityRelationshipType, error)` — `PUT /api/v1/relationship-types/{tag}`
  - `Delete(ctx, tag string) error` — `DELETE /api/v1/relationship-types/{tag}`

  `CreateRelationshipTypeRequest` includes all non-computed fields.
  `UpdateRelationshipTypeRequest` excludes `tag` (path param) but includes all mutable fields.

- [ ] **Step 3: Register client method on `HttpClient`**

  In `internal/cortex/http_client.go`:
  ```go
  func (c *HttpClient) RelationshipTypes() RelationshipTypesClientInterface {
      return &RelationshipTypesClient{client: c}
  }
  ```

- [ ] **Step 4: Build to confirm no compile errors**

  ```bash
  go build ./...
  ```
  Expected: clean.

- [ ] **Step 5: Commit**

  ```bash
  git add internal/cortex/relationship_types.go internal/cortex/http_client.go
  git commit -m "feat: add relationship types API client"
  ```

---

## Task 3: `cortex_relationship_type` resource, data source, and tests

**Files:**
- Create: `internal/provider/relationship_type_models.go`
- Create: `internal/provider/relationship_type_resource.go`
- Create: `internal/provider/relationship_type_data_source_models.go`
- Create: `internal/provider/relationship_type_data_source.go`
- Modify: `internal/provider/provider.go`
- Create: `examples/resources/cortex_relationship_type/resource.tf`
- Create: `examples/data-sources/cortex_relationship_type/data-source.tf`
- Create: `internal/provider/relationship_type_resource_test.go`

**Interfaces:**
- Consumes: `EntityRelationshipType`, `EntityTypeFilter`, `RelationshipInheritance`, `HttpClient.RelationshipTypes()` from Task 2.

- [ ] **Step 1: Write `relationship_type_models.go`**

  ```go
  type RelationshipTypeFilterModel struct {
      Include   types.Bool `tfsdk:"include"`
      Types     types.Set  `tfsdk:"types"`      // SetAttribute(String)
      Providers types.Set  `tfsdk:"providers"`  // SetAttribute(String)
  }

  type RelationshipTypeInheritanceModel struct {
      ParameterTag    types.String `tfsdk:"parameter_tag"`
      InheritanceType types.String `tfsdk:"inheritance_type"`
  }

  type RelationshipTypeResourceModel struct {
      Tag                      types.String `tfsdk:"tag"`
      Name                     types.String `tfsdk:"name"`
      Description              types.String `tfsdk:"description"`
      AllowCycles              types.Bool   `tfsdk:"allow_cycles"`
      IsSingleSource           types.Bool   `tfsdk:"is_single_source"`
      IsSingleDestination      types.Bool   `tfsdk:"is_single_destination"`
      DefinitionLocation       types.String `tfsdk:"definition_location"`
      SourcesFilter            *RelationshipTypeFilterModel `tfsdk:"sources_filter"`
      DestinationsFilter       *RelationshipTypeFilterModel `tfsdk:"destinations_filter"`
      SourceLabelSingular      types.String `tfsdk:"source_label_singular"`
      SourceLabelPlural        types.String `tfsdk:"source_label_plural"`
      DestinationLabelSingular types.String `tfsdk:"destination_label_singular"`
      DestinationLabelPlural   types.String `tfsdk:"destination_label_plural"`
      Inheritances             []RelationshipTypeInheritanceModel `tfsdk:"inheritances"`
      IsCortexManaged          types.Bool   `tfsdk:"is_cortex_managed"`
  }
  ```

  `FromApiModel(ctx, diags, entity)` — maps all fields; `Types` and `Providers` in filters use `types.SetValueMust(types.StringType, ...)`.
  `ToApiModel()` — maps to `CreateRelationshipTypeRequest`; omit empty string fields.

- [ ] **Step 2: Write `relationship_type_resource.go`**

  Schema follows the spec table. Key attributes:
  - `tag`: Required, `stringplanmodifier.RequiresReplace()` (tag is immutable; changing it requires a new resource)
  - `is_cortex_managed`: Computed, `boolplanmodifier.UseStateForUnknown()`
  - `sources_filter` / `destinations_filter`: `SingleNestedAttribute`, Required
  - `sources_filter.types` / `providers`: `SetAttribute(types.StringType)`, Optional
  - `destinations_filter.types` / `providers`: `SetAttribute(types.StringType)`, Optional
  - `inheritances`: `ListNestedAttribute`, Optional

  `ImportState`: call `RelationshipTypes().Get(ctx, req.ID)`. If `entity.IsCortexManaged == true`, add error:
  ```
  "Cannot Import Cortex-Managed Relationship Type"
  "<tag> is managed by Cortex and cannot be updated or deleted via the API. Use terraform state rm if needed."
  ```
  Otherwise, `resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)`.

- [ ] **Step 3: Write `relationship_type_data_source_models.go` and `relationship_type_data_source.go`**

  Mirror the resource model with all fields Computed. Data source looks up by `tag` (Required).

- [ ] **Step 4: Register in `provider.go`**

  Add `NewRelationshipTypeResource` to `Resources()` and `NewRelationshipTypeDataSource` to `DataSources()`.

- [ ] **Step 5: Build to confirm clean compile**

  ```bash
  go build ./...
  ```

- [ ] **Step 6: Write failing acceptance tests**

  `TestAccRelationshipTypeResourceMinimal`:
  - Create with required fields only (no labels, no inheritances, `types = []string{}` for filters)
  - Check `tag`, `name`, `allow_cycles`, `is_single_source`, `is_single_destination`
  - ImportState step
  - Update step (change `name`); verify no phantom diff

  `TestAccRelationshipTypeResourceFull`:
  - Create with labels, `inheritances = [{parameter_tag = "x-cortex-owners", inheritance_type = "FALLBACK"}]`
  - Check all optional label fields
  - Verify plan is clean after apply (no diff from unordered set fields)

  Use unique tags like `"test-reltype-minimal"` and `"test-reltype-full"`.

  Run:
  ```bash
  go test -v ./internal/provider -run TestAccRelationshipType
  ```
  Expected: FAIL — resource type not registered yet (before step 4 is complete).

- [ ] **Step 7: Run acceptance tests to verify they pass**

  ```bash
  go test -v ./internal/provider -run TestAccRelationshipType
  ```
  Expected: PASS.

- [ ] **Step 8: Write example files**

  `examples/resources/cortex_relationship_type/resource.tf`:
  ```hcl
  resource "cortex_relationship_type" "example" {
    tag                   = "service-to-database"
    name                  = "Service to Database"
    description           = "Tracks which services depend on which databases."
    allow_cycles          = false
    is_single_source      = false
    is_single_destination = false
    definition_location   = "SOURCE"

    sources_filter = {
      include = true
      types   = ["service"]
    }

    destinations_filter = {
      include = true
      types   = ["resource"]
    }

    source_label_singular = "Service"
    source_label_plural   = "Services"

    inheritances = [
      {
        parameter_tag    = "x-cortex-owners"
        inheritance_type = "FALLBACK"
      }
    ]
  }
  ```

  `examples/data-sources/cortex_relationship_type/data-source.tf`:
  ```hcl
  data "cortex_relationship_type" "example" {
    tag = "service-to-database"
  }
  ```

- [ ] **Step 9: Run `go generate`, lint, format**

  ```bash
  go generate ./...
  gofmt -w ./...
  golangci-lint run ./...
  ```

- [ ] **Step 10: Commit**

  ```bash
  git add internal/provider/relationship_type_*.go \
          internal/provider/provider.go \
          examples/resources/cortex_relationship_type/ \
          examples/data-sources/cortex_relationship_type/ \
          internal/provider/relationship_type_resource_test.go \
          docs/
  git commit -m "feat: add cortex_relationship_type resource and data source"
  ```

---

## Task 4: Build provider locally and wire up MLB test project

**Files:**
- Create/modify: `~/.terraformrc`
- Modify: `/Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb/main.tf`

**Interfaces:**
- Consumes: compiled provider binary from Task 3.

- [ ] **Step 1: Build provider locally**

  ```bash
  cd /Users/jeffschnitter/git/terraform-provider-cortex
  go build -o ~/.terraform.d/plugins/registry.terraform.io/cortexapps/cortex/99.0.0/darwin_arm64/terraform-provider-cortex .
  ```

- [ ] **Step 2: Create `~/.terraformrc` with dev_overrides**

  ```hcl
  provider_installation {
    dev_overrides {
      "cortexapps/cortex" = "/Users/jeffschnitter/.terraform.d/plugins/registry.terraform.io/cortexapps/cortex/99.0.0/darwin_arm64"
    }
    direct {}
  }
  ```

- [ ] **Step 3: Switch MLB project to local backend**

  In `/Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb/main.tf`, comment out the `cloud {}` block and add:
  ```hcl
  backend "local" {}
  ```
  (The existing `terraform.tfstate` already provides local state.)
  Keep `version = "~> 0.1"` — dev_overrides bypasses version checking.

- [ ] **Step 4: Verify provider loads**

  ```bash
  cd /Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb
  terraform providers
  ```
  Expected: shows `cortexapps/cortex` from the dev override path, no version error.

- [ ] **Step 5: Commit MLB `main.tf` change**

  ```bash
  cd /Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb
  git add main.tf
  git commit -m "chore: use local provider build for entity/relationship type testing"
  ```

---

## Task 5: Add MLB test data and run plan + apply

**Files:**
- Modify: `/Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb/resource_definition.tf`
- Create: `/Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb/relationship_type.tf`
- Create: `/Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb/leaders.tf`

- [ ] **Step 1: Add `mlb-leader` entity type to `resource_definition.tf`**

  ```hcl
  resource "cortex_resource_definition" "mlb_leader" {
    type        = "mlb-leader"
    name        = "MLB Statistical Leader"
    description = "Top 10 leaders for an MLB statistical category in a given season."
    schema = jsonencode({
      type = "object"
      properties = {
        category = { type = "string", description = "Statistical category (HR, RBI, AVG, Wins, ERA, K)" }
        season   = { type = "integer", description = "MLB season year" }
        league   = { type = "string", enum = ["AL", "NL"] }
        leaders = {
          type = "array"
          items = {
            type = "object"
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

- [ ] **Step 2: Create `relationship_type.tf`**

  ```hcl
  resource "cortex_relationship_type" "mlb_to_leader" {
    tag                   = "mlb-to-leader"
    name                  = "MLB Season to Statistical Leader"
    description           = "Links an MLB season entity to its statistical leader boards."
    allow_cycles          = false
    is_single_source      = false
    is_single_destination = false
    definition_location   = "SOURCE"

    sources_filter = {
      include = true
      types   = ["mlb"]
    }

    destinations_filter = {
      include = true
      types   = ["mlb-leader"]
    }

    inheritances = []

    depends_on = [cortex_resource_definition.mlb_leader]
  }
  ```

- [ ] **Step 3: Create `leaders.tf` with 12 leader entities**

  Create `cortex_catalog_entity` resources of type `mlb-leader` for:
  AL: `2026-al-leader-hr`, `2026-al-leader-rbi`, `2026-al-leader-avg`, `2026-al-leader-wins`, `2026-al-leader-era`, `2026-al-leader-strikeouts`
  NL: same six with `nl`.

  Each entity has `tag`, `name`, `type = "mlb-leader"`, and a `description`. Look at the existing entity pattern in the 2026 season module for the exact `cortex_catalog_entity` schema.

- [ ] **Step 4: Run `terraform plan`**

  ```bash
  cd /Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb
  terraform plan
  ```
  Expected: plan shows creates for `mlb-leader` entity type, `mlb-to-leader` relationship type, 12 leader entities. No unexpected diffs on existing resources.

- [ ] **Step 5: Run `terraform apply`**

  ```bash
  terraform apply
  ```
  Expected: all resources created successfully, exit 0.

- [ ] **Step 6: Run `terraform plan` again**

  ```bash
  terraform plan
  ```
  Expected: `No changes. Infrastructure is up-to-date.` (verifies no phantom diffs from unordered sets or missing `UseStateForUnknown`).

- [ ] **Step 7: Commit MLB additions**

  ```bash
  cd /Users/jeffschnitter/git/jeff-test-org/cortex-terraform-mlb
  git add resource_definition.tf relationship_type.tf leaders.tf
  git commit -m "feat: add mlb-leader entity type, relationship type, and 2026 leader entities"
  ```

---

## Task 6: Final cleanup for provider repo

- [ ] **Step 1: Run full unit test suite**

  ```bash
  cd /Users/jeffschnitter/git/terraform-provider-cortex
  make test
  ```
  Expected: all pass.

- [ ] **Step 2: Final lint and format pass**

  ```bash
  gofmt -w ./...
  golangci-lint run ./...
  ```
  Expected: clean.

- [ ] **Step 3: Confirm docs are current**

  ```bash
  go generate ./...
  git diff docs/
  ```
  Expected: no uncommitted doc changes. If diff, commit them.

- [ ] **Step 4: Final commit if needed**

  ```bash
  git add docs/
  git commit -m "docs: regenerate provider docs"
  ```
