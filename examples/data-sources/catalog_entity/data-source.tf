data "cortex_catalog_entity" "products-service" {
  tag = "products-service"
}

output "products_service_name" {
  value = data.cortex_catalog_entity.products-service.name
}

output "products_service_type" {
  value = data.cortex_catalog_entity.products-service.type
}

output "products_service_groups" {
  value = data.cortex_catalog_entity.products-service.groups
}

output "products_service_owners" {
  value = data.cortex_catalog_entity.products-service.ownership
}
