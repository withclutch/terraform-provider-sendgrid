#!/bin/bash

# Import an existing event webhook using its ID
terraform import sendgrid_event_webhook.main <webhook-id>

# To find your webhook IDs, use the SendGrid API:
# curl -X GET "https://api.sendgrid.com/v3/user/webhooks/event/settings/all" \
#   -H "Authorization: Bearer $SENDGRID_API_KEY"
