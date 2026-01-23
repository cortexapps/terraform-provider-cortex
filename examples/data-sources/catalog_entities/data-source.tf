# Search for all service entities
data "cortex_catalog_entities" "all_services" {
  types = ["service"]
}

# Search for entities by query text
data "cortex_catalog_entities" "search" {
  query = "backend"
}

# Search for entities owned by a specific team
data "cortex_catalog_entities" "team_entities" {
  owners = ["backend-team"]
  types  = ["service"]
}

# Search for entities in specific groups
data "cortex_catalog_entities" "production_services" {
  groups = ["production"]
  types  = ["service"]
}

# Search for entities by git repository
data "cortex_catalog_entities" "repo_entities" {
  git_repositories = ["myorg/my-repo"]
}

# Include archived entities in results
data "cortex_catalog_entities" "all_including_archived" {
  include_archived = true
}

# Access the list of entities
output "entity_tags" {
  value = [for entity in data.cortex_catalog_entities.all_services.entities : entity.tag]
}

output "entity_names" {
  value = [for entity in data.cortex_catalog_entities.all_services.entities : entity.name]
}
