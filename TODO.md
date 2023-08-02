
General TODO file of things left to implement.

### General

- ImportState is unverified and likely doesn't work in most elements
- Add more acceptance tests for various changing of elements

### Catalog Entity

Entities:
- x-cortex-team type of Entity (new feature, see Teams section) 
- x-cortex-infra (AWS, GCP, etc)
- x-cortex-k8s
- OpenAPI docs?
- Packages?
- `x-cortex-domain-parents`:
    tag: domain-iam

### Scorecards

Needs more complete acceptance tests. Only basic tests for now.

### Teams

Redo all of this to support the new change of Teams to Entities. This likely means moving to YAML based API approach
for Teams.

### Departments

Needs more full end-to-end tests.

### Resource Definitions

No support yet.

