# Multi-instance integrations: <integration>/<alias>
terraform import cortex_integration_configuration.datadog datadog/datadog-prod

# Single-instance integrations have one configuration per tenant: <integration>
terraform import cortex_integration_configuration.pagerduty pagerduty

# The Cortex API never returns secrets, so set credentials in the configuration. Also set every setting to its value
# in Cortex. A setting that differs or is left out (for example a Datadog custom subdomain) changes Cortex on the
# next apply. A setting that the API cannot change in place (for example the Jira host, or a GitLab host that is left
# out) replaces the configuration.
