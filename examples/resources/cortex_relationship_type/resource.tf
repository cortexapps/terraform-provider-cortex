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
