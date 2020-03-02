resource "sysdig_secure_rule_container" "sample" {
  name        = "Other example of Policy"
  description = "this is other example of policy"
  tags        = ["container", "cis"]

  matching   = true // default
  containers = ["foo", "foo:bar"]
}

resource "sysdig_secure_rule_filesystem" "foo" {
  name        = "Other example of Policy"
  description = "this is other example of policy"
  tags        = ["filesystem", "cis"]

  read_only {
    matching = true // default
    paths    = ["/etc"]
  }

  read_write {
    matching = true // default
    paths    = ["/tmp"]
  }
}

resource "sysdig_secure_rule_network" "foo" {
  name        = "Other example of Policy" // ID
  description = "this is other example of policy"
  tags        = ["network", "cis"]

  block_inbound  = true
  block_outbound = true

  tcp {
    matching = true // default
    ports    = [80, 443]
  }

  udp {
    matching = true // default
    ports    = [80, 443]
  }
}

resource "sysdig_secure_rule_process" "foo" {
  name        = "Other example of Policy" // ID
  description = "this is other example of policy"

  matching  = true // default
  processes = ["bash"]
}

resource "sysdig_secure_rule_syscall" "foo" {
  name        = "Other example of Policy" // ID
  description = "this is other example of policy"

  matching = true // default
  syscalls = ["open", "execve"]
}

resource "sysdig_secure_rule_falco" "foo" {
  name        = "Other example of Policy" // ID
  description = "this is other example of policy"
  tags        = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output    = "A shell was spawned in a container with an attached terminal (user=%user.name %container.info shell=%proc.name parent=%proc.pname cmdline=%proc.cmdline terminal=%proc.tty container_id=%container.id image=%container.image.repository)"
  priority  = "notice"
  source    = "syscall" // syscall or k8s_audit
}
