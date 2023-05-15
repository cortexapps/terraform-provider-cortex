resource "cortex_team" "platform-engineering" {
  tag         = "platform-engineering"
  name        = "Platform Engineering"
  description = "The platform engineering team"
  summary     = "This represents the Cortex platform engineering team"

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
