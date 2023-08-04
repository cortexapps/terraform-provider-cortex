Changelog for the Cortex terraform provider.

## Unreleased

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
