<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/.github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>


Terraform Provider for Sysdig
=============================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
- [Blog on how to use this provider with Sysdig Secure](https://sysdig.com/blog/using-terraform-for-container-security-as-code/)


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) > 0.12.x
-	[Go](https://golang.org/doc/install) > 1.15 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/draios/terraform-provider-sysdig`

```sh
$ git clone git@github.com:draios/terraform-provider-sysdig
$ cd terraform-provider-sysdig
$ make build
```

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.


Contribute
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-sysdig
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

If you want to execute the acceptance tests, you can run `make testacc`.
```sh
$ make testacc
```

<br/>:warning:Please note that you need a token for Monitor and Secure, and since the **acceptance tests create real infrastructure**
you should execute them in an environment where you can remove the resorces easily.



### Creating new resource / data sources

TL;DR;
- Create the resource/data source item
- Add the created item into the `provider.go` resource or datasource map with its wiring
- With its [acceptance test](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests)
- Add its documentation page on `./website/docs/`

### Proposing PR's

* on pull-requests some validations are enforced.
  this can be prevented using [**pre-commit**](https://pre-commit.com)
  * Defined in [`/.pre-commit-config.yaml`](https://github.com/sysdiglabs/terraform-provider-sysdig/blob/master/.pre-commit-config.yaml)
* for `testacc` some credentials are required, check [`/.envrc.template`](https://github.com/sysdiglabs/terraform-provider-sysdig/blob/master/.envrc.template)


### Release

* Use **semver** for releases https://semver.org
* To create a new release, create and push a new tag and it will be released  by [`/.github/workflows/release.yml`](https://github.com/sysdiglabs/terraform-provider-sysdig/blob/master/.github/workflows/release.yml)
