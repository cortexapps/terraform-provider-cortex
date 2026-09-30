resource "cortex_catalog" "example" {
  slug     = "my-services"
  name     = "My Services"
  icon_tag = "cortex"
  is_draft = false

  description = "All microservices in production"
  type        = "FILTER"

  filter {
    query = "tag != null"

    types {
      include = ["service"]
    }

    groups {
      include = ["platform", "backend"]
    }
  }
}
