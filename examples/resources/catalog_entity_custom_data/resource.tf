resource "cortex_catalog_entity_custom_data" "data-with-single-key" {
  tag         = "products-service"
  key         = "region"
  description = "The region where the product is available"
  value       = "us-central1"
}

resource "cortex_catalog_entity_custom_data" "data-with-nested-key" {
  tag         = "products-service"
  key         = "deployment pipeline"
  description = "Deployment pipeline information"
  value       = jsonencode({
    "region" : "us-central1",
    "environments" : [
      "integration",
      "staging",
      "production",
    ],
    "instances" : 3,
  })
}
