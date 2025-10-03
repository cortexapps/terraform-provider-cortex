# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Terraform provider for Cortex.io (https://www.cortexapp.com), built using the terraform-plugin-framework. It allows users to manage Cortex resources (catalog entities, departments, scorecards, resource definitions, etc.) through Terraform.

## Development Commands

### Environment Setup
- Copy `.env-example` to `.env` and set `CORTEX_API_TOKEN` before running tests
- The Makefile loads environment variables from `.env` file automatically

### Build and Installation
- `make build` - Build the provider binary to `./bin/terraform-provider-cortex`
- `make install` - Build and install to local Terraform plugins directory (`~/.terraform.d/plugins/`)
- `make release` - Build release binaries for Linux (amd64) and Darwin (amd64, arm64)

### Testing
- `make test` - Run unit tests (clears test cache first)
- `make testacc` - Run acceptance tests (requires valid `CORTEX_API_TOKEN`, creates real resources)
- `go test -v ./internal/provider -run TestName` - Run a specific test

### Code Quality
- `make lint` - Run golangci-lint
- `make format` - Format code with go fmt
- `make docs` - Generate provider documentation via `go generate`

### Local Development
To use a locally built provider:
1. Build: `go build -o build/terraform-provider-cortex .`
2. Create `~/.terraformrc` with dev_overrides pointing to the build directory
3. Use `source = "cortexlocal/cortex"` in your terraform configuration

## Architecture

### Package Structure
- `main.go` - Provider entry point, runs the provider server
- `internal/provider/` - Terraform provider implementation (resources, data sources, schemas)
- `internal/cortex/` - Cortex API client library

### API Client Architecture (`internal/cortex/`)
The HTTP client uses a functional options pattern for initialization:
- `HttpClient` - Core client with two sling clients (JSON and YAML decoders)
- Client interfaces accessed via methods like `client.CatalogEntities()`, `client.Teams()`, etc.
- `BaseUris` map defines API endpoints for different resource types
- `Route(domain, path)` helper constructs full API paths
- Error handling via `ApiError` type with special handling for 404 and 401

### Provider Architecture (`internal/provider/`)
- `provider.go` - Main provider configuration, registers all resources and data sources
- Resources implement `resource.Resource` and `resource.ResourceWithImportState` interfaces
- Data sources implement `datasource.DataSource` interface
- Each resource/data source has corresponding `*_models.go` files for Terraform state models
- Models use `tfsdk` struct tags to map to Terraform schema attributes

### Resource/Data Source Pattern
Each resource follows this structure:
1. Resource struct with `client *cortex.HttpClient` field
2. Model struct(s) with `types.String`, `types.List`, etc. for state management
3. CRUD methods: `Create`, `Read`, `Update`, `Delete` (resources only)
4. Schema definition with attributes, validators, and plan modifiers
5. Conversion methods between Terraform models and API models (e.g., `ToApiModel`, `FromApiResponse`)

### Available Resources
- `cortex_catalog_entity` - Catalog entities (services, resources, teams, domains)
- `cortex_catalog_entity_custom_data` - Custom data for catalog entities
- `cortex_catalog_entity_openapi` - OpenAPI specs for catalog entities (YAML format)
- `cortex_department` - Departments
- `cortex_resource_definition` - Resource type definitions
- `cortex_scorecard` - Scorecards

### Available Data Sources
- `cortex_catalog_entity` - Read catalog entities
- `cortex_catalog_entity_custom_data` - Read custom data
- `cortex_department` - Read departments
- `cortex_resource_definition` - Read resource definitions
- `cortex_scorecard` - Read scorecards
- `cortex_team` - Read teams

### Special Parsers
- `catalog_entity_parser.go` - Handles conversion between Cortex API entity format and Terraform state
- `scorecard_parser.go` - Handles scorecard-specific conversions and validation

## Provider Configuration

Environment variables:
- `CORTEX_API_TOKEN` - Required API token for authentication
- `CORTEX_API_URL` - Optional API URL (defaults to `https://api.getcortexapp.com`)
- `HTTP_DEBUG=1` - Enable HTTP request/response logging via go-loghttp

## Important Notes

- The provider uses version from Makefile (`VERSION=0.4.6-dev`), injected at build time via ldflags
- Acceptance tests create real resources and require a valid API token
- Documentation is auto-generated via terraform-plugin-docs (run `go generate`)
- The `catalog_entity_openapi` resource uses YAML decoder for OpenAPI specs
- Many resources support both JSON and YAML format responses from the API
