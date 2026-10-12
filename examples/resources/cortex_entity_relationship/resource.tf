resource "cortex_entity_relationship" "example" {
  source_entity_tag     = "products-service"
  relationship_type_tag = "service-to-database"
  destination_entity_tags = [
    "products-db",
    "variants-db",
  ]
}
