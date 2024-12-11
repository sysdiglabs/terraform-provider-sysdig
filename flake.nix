{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    let
      overlays.default = final: prev: {
        terraform-providers = prev.terraform-providers // {
          sysdig = prev.callPackage ./package.nix { };
        };
      };
      flake = flake-utils.lib.eachDefaultSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
            overlays = [ self.overlays.default ];
          };
        in
        {
          packages = with pkgs.terraform-providers; {
            inherit sysdig;
            default = sysdig;
          };
          apps.terraform = flake-utils.lib.mkApp {
            drv = pkgs.terraform.withPlugins (tf: [ tf.sysdig ]);
          };
          devShells.default =
            with pkgs;
            mkShell {
              packages = [
                go_1_23
                govulncheck
              ];
            };

          formatter = pkgs.nixfmt-rfc-style;
        }
      );
    in
    flake // { inherit overlays; };
}
