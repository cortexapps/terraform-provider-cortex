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

resource "cortex_integration_configuration" "gitlab" {
  alias       = "gitlab"
  credentials = { token = { value = var.gitlab_token } }
  gitlab = {
    host        = "https://gitlab.acme.internal"
    group_names = ["platform"]
  }
}

resource "cortex_integration_configuration" "incident_io" {
  alias       = "incident-io"
  credentials = { token = { value = var.incident_io_api_key } }
  incident_io = {}
}

resource "cortex_integration_configuration" "jira" {
  alias = "jira-cloud"
  credentials = {
    basic = { username = "bot@acme.com", password = var.jira_api_token }
  }
  jira = {
    cloud = { subdomain = "acme", base_url = "atlassian.net" }
  }
}

resource "cortex_integration_configuration" "pagerduty" {
  credentials = { token = { value = var.pagerduty_token } }
  pagerduty   = { is_token_readonly = true }
}
