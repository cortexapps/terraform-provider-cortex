terraform {
  required_providers {
    cortex = {
      source = "cortexlocal/cortex"
    }
  }
}

provider "cortex" {
  token = "access-token-here"
}

resource "cortex_catalog_entity_openapi" "test_service_oas1" {
  entity_tag = "test-service"
  spec       = <<EOF
{
    "openapi": "3.0.0",
    "info": {
        "title": "test-service",
        "version": "1.0.0"
    },
    "paths": {
        "/test": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    }
}
EOF
}