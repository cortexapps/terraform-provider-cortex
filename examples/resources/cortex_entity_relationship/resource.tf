resource "cortex_entity_relationship" "example" {
  source_entity_tag      = "al-2026"
  relationship_type_tag  = "mlb-to-leader"
  destination_entity_tags = [
    "2026-al-leader-hr",
    "2026-al-leader-rbi",
  ]
}
