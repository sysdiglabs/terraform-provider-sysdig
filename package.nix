{
  terraform-providers,
}:
# Allows testing of the provider with the following block:
# terraform {
#   required_providers {
#     sysdig = {
#       source  = "sysdiglabs/sysdig"
#       version = "=1.0.0-local"
#     }
#   }
# }
terraform-providers.mkProvider {
  owner = "sysdiglabs";
  repo = "terraform-provider-sysdig";
  homepage = "https://registry.terraform.io/providers/sysdiglabs/sysdig";
  rev = "1.0.0-local"; # Keeping this version fixed with a `-local` version, so user can just bundle the concrete plugin version with terraform using nix.
  vendorHash = "sha256-L+XwC7c4ph4lM0+BhHB9oi1R/Av8jlDcqHewOmtPU1s=";
  hash = "";
  mkProviderFetcher = { ... }: ./.;
}
