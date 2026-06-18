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

# Include ownership information for each entity
data "cortex_catalog_entities" "services_with_owners" {
  types          = ["service"]
  include_owners = true
}

# Access the list of entities
output "entity_tags" {
  value = [for entity in data.cortex_catalog_entities.all_services.entities : entity.tag]
}

output "entity_names" {
  value = [for entity in data.cortex_catalog_entities.all_services.entities : entity.name]
}

# Access git repository for each service (null if no git integration is configured)
output "entity_git_repositories" {
  value = {
    for entity in data.cortex_catalog_entities.all_services.entities :
    entity.tag => entity.git.repository if entity.git != null
  }
}

# Access the first owning team tag for each service (requires include_owners = true)
output "entity_owners" {
  value = {
    for entity in data.cortex_catalog_entities.services_with_owners.entities :
    entity.tag => entity.ownership.groups[0].group_name
    if entity.ownership != null && length(entity.ownership.groups) > 0
  }
}
