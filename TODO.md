
General TODO file of things left to implement.

### General

- Verify we're actually writing data to state file
- ImportState is unverified and likely doesn't work
- Add more acceptance tests for various changing of elements
- Ensure generated YAML in Cortex is valid, specifically with regards to empty attributes

### Catalog Entity

Entities:
- x-cortex-infra (AWS, GCP, etc)
- x-cortex-k8s
- OpenAPI docs?
- Packages?

### Scorecards

Connect client to TF provider. Undo YAML approach (work with Cortex to figure out spec.)

### Teams

Fix the UPDATE issue (this may be a bug with Cortex API).

Also:
- Team Hierarchies? Department mappings to teams?

### Resource Definitions

Get provider working, figure out acceptance tests.

### All Data Sources

Connect to Cortex TF account, figure out how to test safely.
