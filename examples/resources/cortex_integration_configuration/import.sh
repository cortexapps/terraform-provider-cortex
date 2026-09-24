# Multi-instance integrations: <integration>/<alias>
terraform import cortex_integration_configuration.datadog datadog/datadog-prod

# The Cortex API never returns secrets, so set credentials in the configuration. Also set settings that the API
# cannot change in place (for example the Datadog region) to the values in Cortex, or the next apply replaces the
# configuration.
