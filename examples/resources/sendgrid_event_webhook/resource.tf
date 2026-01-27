# Basic event webhook configuration
resource "sendgrid_event_webhook" "basic" {
  enabled       = true
  url           = "https://api.myapp.com/sendgrid/events"
  friendly_name = "Production Event Webhook"

  # Basic email events
  delivered   = true
  bounce      = true
  dropped     = true
  spam_report = true
  unsubscribe = true

  # Engagement events (requires tracking enabled)
  open  = true
  click = true

  # Processing events
  processed = true
  deferred  = true

  # Group events (requires subscription tracking)
  group_resubscribe = true
  group_unsubscribe = true
}

# Multiple webhooks example
resource "sendgrid_event_webhook" "dev" {
  enabled       = true
  url           = "https://dev.myapp.com/sendgrid/events"
  friendly_name = "Development Webhook"

  delivered = true
  bounce    = true
  dropped   = true
}

resource "sendgrid_event_webhook" "prod" {
  enabled       = true
  url           = "https://prod.myapp.com/sendgrid/events"
  friendly_name = "Production Webhook"

  delivered   = true
  bounce      = true
  dropped     = true
  spam_report = true
  open        = true
  click       = true
}
