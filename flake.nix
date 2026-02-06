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
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in
      {
        devShells.default =
          with pkgs;
          mkShell {
            packages = [
              go
              terraform
              goreleaser
              gnupg
              golangci-lint
              gofumpt
              jq
              gnumake
              pre-commit
            ];

            shellHook = ''
              pre-commit install
            '';
          };

        formatter = pkgs.nixfmt-tree;
      }
    );
}
