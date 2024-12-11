{
  terraform-providers,
}:
terraform-providers.mkProvider {
  owner = "sysdiglabs";
  repo = "terraform-provider-sysdig";
  homepage = "https://registry.terraform.io/providers/sysdiglabs/sysdig";
  rev = "master";
  vendorHash = "sha256-9ru4RkH2fDWcgM0I3URlWd811PwySktd+gLsEr624WM=";
  hash = "";
  mkProviderFetcher = { ... }: ./.;
}
