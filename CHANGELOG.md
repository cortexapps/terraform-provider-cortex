Changelog for the Cortex terraform provider.

## Unreleased

## 0.1.11

* Handle null description values on catalog entity dependencies 
* Adjustments to NewRelic APM to allow for multiple applications
* Bump `life4/genesis` from 1.4 to 1.7
* Better handling of empty values of Description fields across models
* Catalog Entities: Fast fail during errors if issue in parsing tf configuration occurs

## 0.1.10

* Fix casting issue with int values in prometheus SLOs

## 0.1.9

* Fix detection mechanism for `cortex_catalog_entity_custom_data` and complex data types

## 0.1.8

* Add Buildkite support for Catalog Entities
* Add support for `x-cortex-slack` on catalog entities

## 0.1.7

* Add `domain_parents` support to catalog entities
* Add `cortex_scorecard` data source for retrieving scorecards
* Add `cortex_scorecard` resource for managing cortex scorecards
* Add `cortex_resource_definition` data source for retrieving resource definitions
* Add `cortex_catalog_entity_custom_data` data source for retrieving scoped custom data on catalog entities
* Add `cortex_catalog_entity_custom_data` resource for managing scoped custom data on catalog entities

## 0.1.6

* Full resource definition support

## 0.1.5

* Fix provider null issue when a domain is created with no children

## 0.1.4

* Fix issue with duplicated query params in parallel subsequent requests in cortex client

## 0.1.2

* Fix issue with XMatters null type being omitted in state cases when not used

## 0.1.1

* Omit `defaultJql` from JIRA issue management when empty

## 0.1.0

* Add `children` support for Domain Catalog Entities
* Add `xmatters` support for OnCall for Catalog Entities
* Add `wiz` support for Catalog Entities
* Add `firehydrant` support for Catalog Entities
* Add the ability to ignore metadata on a Catalog Entity, allowing it to be managed via API instead

## 0.0.10

* Better handle empty values for description and type on Catalog Entities

## 0.0.9

* Add ability to manage custom resources
 
## 0.0.8

* Fix issue with empty metadata responses on x-cortex-dependency entities

## 0.0.7

* Fix issue where `sumo_logic` was improperly required for `slos` block

## 0.0.6

* Fix issue with owners and notificationsEnabled parsing in every type unnecessarily

## 0.0.1

* Initial release
