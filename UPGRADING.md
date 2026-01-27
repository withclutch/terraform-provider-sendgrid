# Upgrading Guide

This document describes breaking changes and how to upgrade between major versions of the SendGrid Terraform Provider.

## Upgrading to v2.1.0

### Event Webhook Breaking Changes

Version 2.1.0 introduces multi-webhook support for the `sendgrid_event_webhook` resource. This is a breaking change that affects existing state.

#### What Changed

1. **Multiple webhooks now supported**: You can create multiple event webhooks per account using `for_each` or multiple resource blocks.

2. **Resource ID format changed**: The resource ID is now the actual SendGrid webhook ID instead of `"default"`.

3. **Delete behavior changed**: Previously, `terraform destroy` would NOT delete event webhooks from SendGrid (no-op). Now it WILL delete them.

#### Manual State Migration Required

If you have existing `sendgrid_event_webhook` resources, you must migrate your Terraform state manually.

**Step 1: Backup your state**

```bash
terraform state pull > backup.tfstate
```

**Step 2: Get your webhook ID from SendGrid**

```bash
# List all webhooks to find the ID
curl -X GET "https://api.sendgrid.com/v3/user/webhooks/event/settings/all" \
  -H "Authorization: Bearer $SENDGRID_API_KEY"
```

The response will include each webhook's `id` field.

**Step 3: Remove old state and import with new ID**

```bash
# Remove the old state entry
terraform state rm sendgrid_event_webhook.example

# Import with the actual webhook ID from step 2
terraform import sendgrid_event_webhook.example <webhook-id>
```

**Step 4: Verify**

```bash
terraform plan  # Should show no changes
```

#### Preserving Webhooks on Destroy

If you want to prevent accidental deletion of webhooks (matching the old behavior), add a lifecycle block:

```hcl
resource "sendgrid_event_webhook" "example" {
  enabled = true
  url     = "https://myapp.com/webhook"
  # ...

  lifecycle {
    prevent_destroy = true
  }
}
```

#### New Features in v2.1.0

- **Multiple webhooks**: Create separate webhooks for different environments or use cases
- **Import support**: Import existing webhooks with `terraform import`
- **Proper deletion**: Webhooks are now actually deleted when the resource is destroyed
- **Better error handling**: Clearer error messages for plan limits and URL conflicts

#### Example: Multiple Webhooks

```hcl
resource "sendgrid_event_webhook" "dev" {
  enabled       = true
  url           = "https://dev.myapp.com/webhook"
  friendly_name = "Development"
  delivered     = true
  bounce        = true
}

resource "sendgrid_event_webhook" "prod" {
  enabled       = true
  url           = "https://prod.myapp.com/webhook"
  friendly_name = "Production"
  delivered     = true
  bounce        = true
  spam_report   = true
}
```
