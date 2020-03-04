resource "sysdig_secure_team" "sample" {
  name               = "sample-team"
  description        = "sample"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""
  use_sysdig_capture = false

  user_roles {
    email = "sample@example.com"
    role  = "ROLE_TEAM_STANDARD"
  }

  user_roles {
    email = "sample2@example.com"
    role  = "ROLE_TEAM_EDIT"
  }

}

