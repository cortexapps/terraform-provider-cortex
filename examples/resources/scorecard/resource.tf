resource "cortex_scorecard" "dora-metrics" {
  tag         = "dora"
  name        = "dora"
  description = "DORA metrics"
  is_draft    = false

  level {
    name   = "Gold"
    number = 1
  }

  level {
    name   = "Silver"
    number = 2
  }

  level {
    name   = "Bronze"
    number = 3
  }

  rule {
    title           = "DORA metric 1"
    expression      = "deploys.frequency >= 1"
    number          = 1
    weight          = 1
    level_name      = "Gold"
    failure_message = "Ship fast, ship often"
    description     = "DORA metric 1 description"
  }
}
