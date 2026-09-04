{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    go-overlay.url = "github:purpleclay/go-overlay";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      self,
      nixpkgs,
      go-overlay,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        useLatestGoVersion = final: prev: {
          go_latest = final.go-bin.latestStable;
        };
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
          overlays = [
            go-overlay.overlays.default
            useLatestGoVersion
          ];
        };
      in
      {
        devShells.default =
          with pkgs;
          mkShell {
            packages = [
              go_latest
              govulncheck
              terraform
              tfproviderdocs
              goreleaser
              gnupg
              golangci-lint
              gofumpt
              jq
              just
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
