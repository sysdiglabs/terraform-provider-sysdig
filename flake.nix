{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    tfproviderdocs-init.url = "github:nixos/nixpkgs/pull/366576/head"; # tfproviderdocs is not yet a package in nixpkgs, so this workaronuds it
  };
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      tfproviderdocs-init,
    }:
    let
      overlays.default = final: prev: {
        terraform-providers = prev.terraform-providers // {
          sysdig = prev.callPackage ./package.nix { };
        };
      };
      overlays.tfproviderdocs =
        final: prev:
        let
          pkgs = import tfproviderdocs-init { inherit (prev) system; };
        in
        {
          inherit (pkgs) tfproviderdocs;
        };
      flake = flake-utils.lib.eachDefaultSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
            overlays = [
              self.overlays.default
              self.overlays.tfproviderdocs
            ];
          };
        in
        {
          # Exposes the local plugin to be consumed as a package.
          packages = with pkgs.terraform-providers; {
            inherit sysdig;
            default = sysdig;
          };

          # To be used with `nix run`.
          apps.terraform = flake-utils.lib.mkApp {
            drv = pkgs.terraform.withPlugins (tf: [ tf.sysdig ]);
          };

          # Used for local development. Adds the required dependencies to work in this project.
          devShells.default =
            with pkgs;
            mkShell {
              packages = [
                go_1_23
                govulncheck
                golangci-lint
                golangci-lint-langserver
                (terraform.withPlugins (tf: [ tf.sysdig ]))
                tfproviderdocs
              ];
            };

          # Used with `nix develop <url/path>#terraform-with-plugin`.
          # You can load terraform with the sysdig plugin from a commit, a branch or a tag.
          # For instance:
          # - `nix develop github:sysdiglabs/terraform-provider-sysdig#terraform-with-plugin` will create a local dev shell with the version from the main branch.
          # - `nix develop github:sysdiglabs/terraform-provider-sysdig/branch-name#terraform-with-plugin` with create a local dev shell with the version from the `branch-name` branch code.
          # - `nix develop github:sysdiglabs/terraform-provider-sysdig/v1.2.3#terraform-with-plugin` will create a local dev shell with the version from the tag `v1.2.3` code (note the provided version is just an example).
          # - `nix develop .#terraform-with-plugin` will create a local dev shell with terraform with the local code.
          devShells.terraform-with-plugin =
            with pkgs;
            mkShell {
              packages = [
                (terraform.withPlugins (tf: [ tf.sysdig ]))
                bashInteractive
              ];
            };

          formatter = pkgs.nixfmt-rfc-style;
        }
      );
    in
    flake // { inherit overlays; };
}
