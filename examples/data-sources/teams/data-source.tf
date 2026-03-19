# Retrieve all teams
data "cortex_teams" "all" {}

# Use team tags to look up individual team details
data "cortex_team" "example" {
  for_each = toset([for t in data.cortex_teams.all.teams : t.tag])
  tag      = each.key
}
