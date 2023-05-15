resource "cortex_team" "engineering" {
  tag         = "engineering"
  name        = "Engineering"
  description = "The engineering team"
  summary     = "This represents the Cortex engineering team"

  link {
    name        = "Twitter"
    description = "Tweet, tweet"
    url         = "https://twitter.com/GetCortexApp"
    type        = "documentation"
  }

  slack_channel {
    name                  = "#engineering"
    notifications_enabled = true
  }

  additional_member {
    name  = "John Doe One"
    email = "test+one@cortex.io"
  }

  additional_member {
    name  = "John Doe Two"
    email = "test+two@cortex.io"
  }
}
