# Cortex Terraform Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.5
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

Add the following block into your Terraform files (note that this provider requires Terraform 1.5 or greater):

```terraform
terraform {
  required_providers {
    cortex = {
      source  = "cortexapps/cortex"
      version = ">= 0.1"
    }
  }
  required_version = ">= 1.5.0"
}
```

You'll need to set the environment variable `CORTEX_API_TOKEN` with your Cortex API token.

You can optionally further configure the provider either inline in terraform:

```terraform
provider "cortex" {
  token        = "my-api-token"
  base_api_url = "https://api.getcortexapp.com" // default value, only set if changing
}
```

...or via ENV:

| Key              | Description                        | Default Value                  |
|------------------|------------------------------------|--------------------------------|
| CORTEX_API_TOKEN | Your Cortex.io API token           | ""                             |
| CORTEX_API_URL   | The base API URL for Cortex's API. | "https://api.getcortexapp.com" |

### Resource Types

This provider comes with the following resource types:

* [`cortex_catalog_entity`](docs/resources/catalog_entity.md)
* [`cortex_catalog_entity_custom_data`](docs/resources/catalog_entity_custom_data.md)
* [`cortex_department`](docs/resources/department.md)
* [`cortex_resource_definition`](docs/resources/resource_definition.md)
* [`cortex_scorecard`](docs/resources/scorecard.md)

And the following data sources:

* [`cortex_catalog_entity`](docs/data-sources/catalog_entity.md)
* [`cortex_catalog_entity_custom_data`](docs/data-sources/catalog_entity_custom_data.md)
* [`cortex_department`](docs/data-sources/department.md)
* [`cortex_resource_definition`](docs/data-sources/resource_definition.md)
* [`cortex_scorecard`](docs/data-sources/scorecard.md)
* [`cortex_team`](docs/data-sources/team.md)

Examples on each of these can be found in the [examples/](examples/) folder.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (
see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin`
directory.

To be able to use the **locally installed provider**, 

1. Run `go build -o build/terraform-provider-cortex .`
This will build the binary and put it in the `build/` directory.

2. Create a `~/.terraformrc` file that lets you override the default provider location:
```terraform
provider_installation {
  dev_overrides {
    "cortexlocal/cortex" = "<path/to/terraform-provider-cortex/build"
  }
  direct {}
}
```

3. Use the `provider` as:

```terraform
terraform {
  required_providers {
    cortex = {
      source  = "cortexlocal/cortex"
    }
  }
}
```

### Documentation

```shell

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.
