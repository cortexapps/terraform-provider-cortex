resource "cortex_integration_configuration" "datadog" {
  alias      = "datadog-prod"
  is_default = true
  credentials = {
    key_pair = { key = var.datadog_api_key, secret = var.datadog_app_key }
  }
  datadog = {
    region       = "US1"
    environments = ["prod"]
  }
}
