# Import by alias. The Cortex API never returns the keys, so set api_key and app_key in the configuration. Also set
# region and custom_subdomain to the values in Cortex: a different value replaces the configuration.
terraform import cortex_datadog_configuration.production datadog-production
