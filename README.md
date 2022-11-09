<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/.github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>


Terraform Provider for Sysdig
=============================

- **[Terraform Registry - Sysdig Provider Docs](https://registry.terraform.io/providers/sysdiglabs/sysdig/latest/docs)**
- [Blog on how to use this provider with Sysdig Secure](https://sysdig.com/blog/using-terraform-for-container-security-as-code/)
- Terraform
  - Website: https://www.terraform.io
  - Mailing list on  [Google Groups](http://groups.google.com/group/terraform-tool)
  - [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)



Contribute
---------------------------


### - Building

#### Requirements

- [Terraform](https://www.terraform.io/downloads.html) > 0.12.x
- [Go](https://golang.org/doc/install) > 1.15 (to build the provider plugin)
  - If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine 
    (version 1.15+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

#### - Developing

First **clone** source repository to: `$GOPATH/src/github.com/draios/terraform-provider-sysdig`

```sh
$ git clone git@github.com:draios/terraform-provider-sysdig
$ cd terraform-provider-sysdig
$ make build
```


To **compile** the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-sysdig
...
```

In order to **test** the provider, you can simply run `make test`.

```sh
$ make test
```

If you want to execute the **acceptance tests**, you can run `make testacc`.

```sh
$ make testacc
```

To use the local provider you just built, follow the instructions to [**install** it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) in your machine with:

```sh
$ make install
```

That will add the provider to the terraform plugins dir. Then just set `source` and `version` values appropriately:

```terraform
provider "aws" {
  region = my_region
}

terraform {
  required_providers {
    sysdig = {
      source = "local/sysdiglabs/sysdig"
      version = "~> 1.0.0"
    }
  }
}
```

To uninstall the plugin:

```sh
$ make uninstall
```

<br/>:warning:Please note that you need a token for Monitor and Secure, and since the **acceptance tests create real infrastructure**
you should execute them in an environment where you can remove the resorces easily.



### - Creating new resource / data sources

TL;DR;
- Create the resource/data source item
- Add the created item into the `provider.go` resource or datasource map with its wiring
- With its [acceptance **test**](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests)
- Add its **documentation** page on `./website/docs/`


Interesting resources
- https://www.terraform.io/plugin
- https://www.hashicorp.com/blog/testing-hashicorp-terraform


### - Proposing PR's

* if it's your first time, validate you're taking into account every aspect of the [`./github/pull_request_template`](.github/pull_request_template.md)
* on pull-requests some validations are enforced.
  - Defined in [`/.pre-commit-config.yaml`](https://github.com/sysdiglabs/terraform-provider-sysdig/blob/master/.pre-commit-config.yaml)
  - You can work on this before even pushing to remote, using [**pre-commit**](https://pre-commit.com) plugin
  
* for Acceptance Tests `testacc` some credentials are required, check [`/.envrc.template`](https://github.com/sysdiglabs/terraform-provider-sysdig/blob/master/.envrc.template)


### -  Release

To create a new release, create and push a new **tag**, and it will be released  following [`/.
github/workflows/release.yml`](https://github.com/sysdiglabs/terraform-provider-sysdig/blob/master/.github/workflows/release.yml).
 
* Before releasing check the **diff** between previous tag and master branch, to spot major changes
* For **tag**, use **[semver](https://semver.org)** 
* Review Released Draft Note, and make it as clear as possible.
* Notify Sysdig teams on our internal #release-announcements slack channel and optionally in #terraform-provider

Mange takk!

![giphy](https://user-images.githubusercontent.com/1073243/200767344-7435f322-24c0-44d2-ac56-468791c84ca5.gif)



