resource "cortex_datadog_configuration" "production" {
  alias        = "datadog-production"
  api_key      = var.datadog_api_key
  app_key      = var.datadog_app_key
  region       = "US1"
  environments = ["prod", "staging"]
  is_default   = true
}
