# Basic domain authentication setup
resource "sendgrid_domain_authentication" "main" {
  domain             = "mycompany.com"
  subdomain          = "em"
  is_default         = true
  automatic_security = true
  custom_spf         = false
}

# Check domain authentication validation
resource "sendgrid_domain_authentication_validation" "this" {
  domain_authentication_id = sendgrid_domain_authentication.main.id
}
