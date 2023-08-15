package provider_test

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityResourceMinimal(t *testing.T) {
	resourceName := "cortex_catalog_entity.test-minimal"
	description := "Minimal configuration service"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceMinimal(description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", "test-minimal"),
					resource.TestCheckResourceAttr(resourceName, "name", "Minimal configuration service"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceMinimal(description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "test-minimal" {
	tag = "test-minimal"
	name = "Minimal configuration service"
	description = "%s"
}
`, description)
}

func TestAccCatalogEntityResourceSimple(t *testing.T) {
	resourceName := "cortex_catalog_entity.test-simple-1"
	description := "Simple configuration service 1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceSimple(description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", "test-simple-1"),
					resource.TestCheckResourceAttr(resourceName, "name", "Simple service"),
					resource.TestCheckResourceAttr(resourceName, "description", description),

					resource.TestCheckResourceAttr(resourceName, "owners.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.name", "John Doe"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.email", "john.doe@cortex.io"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.name", "Engineering"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.provider", "CORTEX"),
					resource.TestCheckResourceAttr(resourceName, "owners.2.type", "SLACK"),
					resource.TestCheckResourceAttr(resourceName, "owners.2.channel", "engineering"),
					resource.TestCheckResourceAttr(resourceName, "owners.2.notifications_enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "groups.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "groups.1", "test2"),

					resource.TestCheckResourceAttr(resourceName, "links.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "links.0.name", "Internal Docs"),
					resource.TestCheckResourceAttr(resourceName, "links.0.type", "documentation"),
					resource.TestCheckResourceAttr(resourceName, "links.0.url", "https://internal-docs.cortex.io/test-simple-1"),

					resource.TestCheckResourceAttr(resourceName, "git.github.repository", "cortexio/test-simple-1"),
					resource.TestCheckResourceAttr(resourceName, "sentry.project", "test-simple-1"),

					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.stream_id", "asdf1234567"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.0.percentile", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.0.target", "2"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.0.slo", "0.9995"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.1.percentile", "0.7"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.1.target", "1"),
					resource.TestCheckResourceAttr(resourceName, "slos.lightstep.streams.0.targets.latencies.1.slo", "0.9998"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceSimple(description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "test-simple-1" {
  tag = "test-simple-1"
  name = "Simple service"
  description = "%s"

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
      url  = "https://internal-docs.cortex.io/test-simple-1"
    }
  ]

  git = {
    github = {
      repository = "cortexio/test-simple-1"
    }
  }

  slos = {
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
  }

  sentry = {
    project = "test-simple-1"
  }
}
`, description)
}

func TestAccCatalogEntityResourceComplete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceComplete("test", "A Test Service", "A test service for the Terraform provider"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "tag", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "name", "A Test Service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "description", "A test service for the Terraform provider"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.#", "4"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.name", "John Doe"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.email", "john.doe@cortex.io"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.type", "GROUP"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.name", "Engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.provider", "CORTEX"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.type", "SLACK"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.channel", "engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.notifications_enabled", "false"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.3.type", "SLACK"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.3.channel", "platform-engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.3.notifications_enabled", "true"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.#", "2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.0", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.1", "test2"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "dependencies.#", "1"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "dependencies.0.tag", "manual-test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "dependencies.0.description", "for testing"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.#", "1"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.name", "Internal Docs"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.type", "documentation"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.url", "https://internal-docs.cortex.io/products-service"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.github.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.gitlab.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.azure.project", "cortexio"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.azure.repository", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "git.bitbucket.repository", "cortexio/products-service"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.0", "123456"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.1", "123457"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids.0", "123456"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids.1", "123457"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.sonar_qube.project", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.veracode.application_names.0", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.veracode.sandboxes.0.application_name", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.veracode.sandboxes.0.sandbox_name", "staging"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "ci_cd.buildkite.pipelines.0.slug", "products-pipeline"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "ci_cd.buildkite.tags.0.tag", "products-tag"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "bug_snag.project", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "sentry.project", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "checkmarx.projects.0.name", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "firehydrant.services.0.id", "asdf1234"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "firehydrant.services.0.type", "ID"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "firehydrant.services.1.id", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "firehydrant.services.1.type", "SLUG"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "rollbar.project", "products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "snyk.projects.0.organization", "cortexio"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "snyk.projects.0.project_id", "cortexio/products-service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "snyk.projects.0.source", "CODE"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "wiz.projects.0.project_id", "01234567-e65f-4b7b-a8b1-5b642894ec37"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_catalog_entity.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityResourceComplete("test", "A Test Service", "A test service for the Terraform provider 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "tag", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "name", "A Test Service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "description", "A test service for the Terraform provider 2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.#", "4"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.#", "2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceComplete(tag string, name string, description string) string {
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
    },
    {
      type                  = "SLACK"
      channel               = "platform-engineering"
      notifications_enabled = true
    }
 ]

  dependencies = [
    { 
      tag = "manual-test"
      description = "for testing"
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
    new_relic = {
      application_id = 123456
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
    }
    gitlab = {
      repository = "cortexio/products-service"
    }
    azure = {
      project    = "cortexio"
      repository = "cortexio/products-service"
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

  ci_cd = {
	buildkite = {
      pipelines = [
        { slug = "products-pipeline" }
	  ]
      tags = [
        { tag = "products-tag" }
      ]
    }
  }

  bug_snag = {
    project = "cortexio/products-service"
  }

  checkmarx = {
    projects = [
      {
        name = "products-service"
      }
    ]
  }

  firehydrant = {
    services = [
      { 
        id   = "asdf1234"
        type = "ID"
      },
      { 
        id   = "products-service"
        type = "SLUG"
      }
    ]
  }

  rollbar = {
    project = "products-service"
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

  wiz = {
	projects = [
	  {
	    project_id = "01234567-e65f-4b7b-a8b1-5b642894ec37"
	  }
	]
  }
}
`, tag, name, description)
}

func TestAccCatalogEntityResourceResourceSimple(t *testing.T) {
	tag := "test-resource-simple"
	resourceName := "cortex_catalog_entity.test-resource-simple"
	name := "Simple Test Resource"
	description := "Simple configuration resource"
	resourceType := "test-resource-definition"
	expectedDefinition := map[string]interface{}{"version": "1.0.0"}
	expectedDefinitionString, _ := json.Marshal(&expectedDefinition)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceResourceSimple(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "type", resourceType),
					resource.TestCheckResourceAttr(resourceName, "definition", string(expectedDefinitionString)),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityResourceResourceSimple(tag, name, "Simple configuration resource 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Simple configuration resource 2"),
					resource.TestCheckResourceAttr(resourceName, "type", resourceType),
					resource.TestCheckResourceAttr(resourceName, "definition", string(expectedDefinitionString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceResourceSimple(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "test-resource-simple" {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 type = "test-resource-definition"
 definition = jsonencode({
	"version": "1.0.0"
 })
}`, tag, name, description)
}

func TestAccCatalogEntityUnmanagedMetadata(t *testing.T) {
	tag := "test-unmanaged-metadata"
	resourceName := "cortex_catalog_entity.test-unmanaged-metadata"
	name := "Unmanaged Metadata Test Entity"
	description := "Unmanaged Metadata entity"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityUnmanagedMetadata(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "ignore_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata",
					"ignore_metadata",
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityUnmanagedMetadata(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "test-unmanaged-metadata" {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 ignore_metadata = true
 metadata = jsonencode({
  "this": "is",
  "ignored": true
 })
}`, tag, name, description)
}
