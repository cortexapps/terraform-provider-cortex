resource "cortex_scorecard" "dora-metrics" {
  tag         = "dora"
  name        = "dora"
  description = "DORA metrics"
  draft       = false

  ladder = {
    levels = [
      {
        name  = "Gold"
        rank  = 3
        color = "#cda400"
      },
      {
        name  = "Silver"
        rank  = 2
        color = "#8c9298"
      },
      {
        name  = "Bronze"
        rank  = 1
        color = "#c38b5f"
      }
    ]
  }
  rules = [
    {
      title           = "DORA metric 1"
      expression      = "deploys.frequency >= 1"
      weight          = 1
      level           = "Gold"
      description     = "DORA metric 1 description"
      failure_message = "Ship fast, ship often"
    }
  ]
  filter = {
    types = {
      include = ["service"]
    }
    query = "entity.description() != null"
  }
  evaluation = {
    window = 24
  }
}
