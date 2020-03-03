resource "sysdig_user" "sample" {
  email      = "sample@example.com"
  system_role = "ROLE_CUSTOMER"
  first_name = "John"
  last_name  = "Smith"
}
