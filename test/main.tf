# test/main.tf file

terraform {
  required_providers {
    sysdig = {
      source  = "local/sysdiglabs/sysdig"
      version = "1.0.0"
    }
  }
}

provider "sysdig" {
  sysdig_secure_api_token = "f866a0da-7bf4-446e-bb15-a1dcd5a2244e"
  sysdig_secure_url       = "https://secure-staging2.sysdig.com"
}

resource "sysdig_secure_posture_policy" "hepl" {
    name = "new policy from terraform 1234"
    description = "new description"
    is_active = true

    groups {
      name = "new group 1"
      description = "new description"  
      groups {
        name = "new group 2"
        description = "new description"  
        requirements {
          name = "req 2"
          description = "new des 2"
          controls {
            name = "/etc/audit/rules.d/*.rules ends with `-e 2` setting"
            enabled = true
          }
        }
        groups {
        name = "new group 3"
        description = "new description"  
        requirements {
          name = "req 3"
          description = "new des 3"
          controls {
            name = "/etc/audit/rules.d/*.rules ends with `-e 2` setting"
            enabled = true
            }
          }
        groups {
          name = "new group 4"
          description = "new description"  
          requirements {
            name = "req 4"
            description = "new des 4"
            controls {
              name = "/etc/audit/rules.d/*.rules ends with `-e 2` setting"
              enabled = true
            }
            }
          groups {
          name = "new group 5 this is the last supported level"
          description = "new description"  
          requirements {
            name = "req 5"
            description = "new des 5"
            controls {
              name = "/etc/audit/rules.d/*.rules ends with `-e 2` setting"
              enabled = true
            }
            }
          }
          }
        }
      }
      requirements {
        name = "req 1"
        description = "new des"
      }
    }

    groups {
      name = "new group a"
      description = "new description"  
    }


}


