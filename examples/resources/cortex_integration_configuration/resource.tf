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

resource "cortex_integration_configuration" "pagerduty" {
  credentials = { token = { value = var.pagerduty_token } }
  pagerduty   = { is_token_readonly = true }
}
