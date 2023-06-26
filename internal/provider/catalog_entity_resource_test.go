package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceConfig("test", "A Test Service", "A test service for the Terraform provider"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "tag", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "name", "A Test Service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "description", "A test service for the Terraform provider"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.#", "3"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.name", "John Doe"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.email", "john.doe@cortex.io"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.type", "GROUP"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.name", "Engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.provider", "CORTEX"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.type", "SLACK"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.channel", "engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.notifications_enabled", "false"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.#", "2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.0", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.1", "test2"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.#", "1"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.name", "Internal Docs"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.type", "documentation"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.url", "https://internal-docs.cortex.io/products-service"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.github.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.github.base_path", "/"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.gitlab.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.gitlab.base_path", "/"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.azure.project", "cortexio"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.azure.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.azure.base_path", "/"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.bitbucket.repository", "cortexio/products-service"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "sentry.project", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "snyk.projects.0.organization", "cortexio"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "snyk.projects.0.project_id", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "snyk.projects.0.source", "CODE"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.code_cov.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.code_cov.provider", "GITHUB"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.0", "123456"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.1", "123457"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids.0", "123456"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids.1", "123457"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.sonar_qube.project", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.veracode.application_names.0", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.veracode.sandboxes.0.application_name", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.veracode.sandboxes.0.sandbox_name", "staging"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_catalog_entity.test",
				ImportState:       true,
				ImportStateVerify: false, // TODO: Fix this
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"tag", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityResourceConfig("test", "A Test Service", "A test service for the Terraform provider 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "tag", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "name", "A Test Service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "description", "A test service for the Terraform provider 2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.#", "3"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.#", "2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceConfig(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "test" {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 
 owners = [
    {
      type  = "EMAIL"
      name  = "John Doe"
      email = "john.doe@cortex.io"
    },
    {
      type     = "GROUP"
      name     = "Engineering"
      provider = "CORTEX"
    },
    {
      type                  = "SLACK"
      channel               = "engineering"
      notifications_enabled = false
    }
 ]

  groups = [
   "test",
   "test2"
  ]

  links = [
    {
      name = "Internal Docs"
      type = "documentation"
      url  = "https://internal-docs.cortex.io/products-service"
    }
  ]

  metadata = jsonencode({
	"my-key": "the value",
	"another-key": {
		"this": "is",
		"an": "object"
	},
	"final-key": [
		"also",
		"use",
		"lists!"
	]
  })

  alerts = [
    {
      type  = "opsgenie"
      tag   = "different-tag"
      value = "my-service-override-tag"
    }
  ]

  apm = {
    data_dog = {
      monitors = [123456, 123457]
    }
    dynatrace = {
      entity_ids           = ["mock-slo-id-1", "mock-slo-id-2"]
      entity_name_matchers = ["products-service", "products-service-2"]
    }
    newrelic = {
      application_id = 123456
      alias          = "products-service"
    }
  }

  dashboards = {
    embeds = [
      {
        type = "grafana"
        url  = "https://grafana.cortex.io/d/123456"
      },
      {
        type = "newrelic"
        url  = "https://newrelic.cortex.io/123456"
      },
      {
        type = "datadog"
        url  = "https://datadog.cortex.io/123456"
      }
    ]
  }

  git = {
    github = {
      repository = "cortexio/products-service"
      base_path  = "/"
    }
    gitlab = {
      repository = "cortexio/products-service"
      base_path  = "/"
    }
    azure = {
      project    = "cortexio"
      repository = "cortexio/products-service"
      base_path  = "/"
    }
    bitbucket = {
      repository = "cortexio/products-service"
    }
  }

  issues = {
    jira = {
      default_jql = "project = CORTEX AND component = Products"
      projects = ["PRODUCTS"]
    }
  }

  slos = {
    data_dog = [
      {
        id = "123456"
      }
    ]
    dynatrace = [
      {
        id = "123456"
      }
    ]
    lightstep = {
      streams = [
        {
          stream_id = "asdf1234567"
          targets = {
            latencies = [
              {
                percentile = 0.5
                target     = 2
                slo        = 0.9995
              },
              {
                percentile = 0.7
                target     = 1
                slo        = 0.9998
              }
            ]
          }
        }
      ]
    }
    prometheus = [
      {
        error_query = "sum(rate(http_requests_total{job=\"products-service\", status=~\"5..\"}[5m])) / sum(rate(http_requests_total{job=\"products-service\"}[5m]))"
        total_query = "sum(rate(http_requests_total{job=\"products-service\"}[5m]))"
        slo         = 0.999
        name        = "HTTP 5xx"
      }
    ]
    signal_fx = [
      {
        query     = "sf_metric:'jvm.memory.max' AND area:'nonheap'"
        rollup    = "AVERAGE"
        target    = 512000
        lookback  = "P1Y"
        operation = "<="
      }
    ]
    sumo_logic = [
      {
        id = "123456"
      }
    ]
  }

  static_analysis = {
    code_cov = {
      repository = "cortexio/products-service"
      provider   = "GITHUB"
    }
    mend = {
      application_ids = ["123456", "123457"]
      project_ids     = ["123456", "123457"]
    }
    sonar_qube = {
      project = "cortexio/products-service"
    }
    veracode = {
      application_names = ["products-service"]
      sandboxes = [
        {
          application_name = "products-service"
          sandbox_name     = "staging"
        }
      ]
    }
  }

  sentry = {
    project = "cortexio/products-service"
  }

  snyk = {
	projects = [
	  {
	    organization = "cortexio"
	    project_id = "cortexio/products-service"
	    source = "CODE"
	  }
	]
  }
}
`, tag, name, description)
}
