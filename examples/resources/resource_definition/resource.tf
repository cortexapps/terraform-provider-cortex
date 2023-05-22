resource "cortex_resource_definition" "squid-proxy" {
  type        = "squid-proxy"
  name        = "Squid Proxy"
  description = "Cortex's customized squid proxy that is used to make requests to firewalled self-managed resources with a static IP."
  schema = {
    properties = {
      ip = {
        type = "string"
      }
    }
  }
}
