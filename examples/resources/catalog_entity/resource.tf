resource "cortex_catalog_entity" "products-service" {
  tag = "products-service"

  name        = "Products Service"
  description = "Serves up products in a nice and hearty fashion."

  owners = [
    {
      name  = "John Doe"
      type  = "user"
      email = "john.doe@cortex.io"
    },
    {
      name        = "Engineering"
      type        = "group"
      description = "The engineering group"
      provider    = "CORTEX"
    },
    {
      type                  = "slack"
      channel               = "engineering"
      notifications_enabled = false
    }
  ]

  groups = [
    "production",
    "lang-golang",
  ]

  links = [
    {
      name = "Internal Docs"
      type = "documentation"
      url  = "https://internal-docs.cortex.io/products-service"
    }
  ]

  metadata = jsonencode({
    "my-key" : "the value",
    "another-key" : {
      "this" : "is",
      "an" : "object"
    },
    "final-key" : [
      "also",
      "use",
      "lists!"
    ]
  })

  dependencies = [
    {
      tag         = "variants-service"
      method      = "POST"
      path        = "/api/v2/variants"
      description = "Creates a new variant"
      metadata = jsonencode({
        "my-key" : "the value",
        "another-key" : {
          "this" : "is",
          "an" : "object"
        },
        "final-key" : [
          "also",
          "use",
          "lists!"
        ]
      })
    }
  ]

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
      entity_ids           = ["123456", "123457"]
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
      default_jql = "project = PRODUCTS AND component = customer-facing"
      projects    = ["PRODUCTS"]
      labels      = ["live"]
      components  = ["customer-facing"]
    }
  }

  on_call = {
    pagerduty = {
      id   = "123456"
      type = "SERVICE" // or SCHEDULE or ESCALATION_POLICY
    }
    ops_genie = {
      id   = "Cortex-Engineering"
      type = "SCHEDULE"
    }
    victor_ops = {
      id   = "team-cortex"
      type = "SCHEDULE"
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
        alias       = "http-5xx"
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
      provider   = "github"
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

  bugsnag = {
    project = "products-service"
  }

  checkmarx = {
    projects = ["products-service"]
  }

  rollbar = {
    project = "products-service"
  }

  sentry = {
    project = "products-service"
  }

  snyk = {
    projects = [
      {
        organization = "cortexio"
        project_id   = "products-service"
      }
    ]
  }
}
