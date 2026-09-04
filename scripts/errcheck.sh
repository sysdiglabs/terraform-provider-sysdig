#!/usr/bin/env bash

# Check gofmt
echo "==> Checking for unchecked errors..."

# nixpkgs' errcheck build lags behind our go.mod version (go/packages refuses to
# load a module declaring a newer 'go' directive than the errcheck binary was
# built with), so it isn't in the flake devShell. Install it on demand instead.
if ! which errcheck > /dev/null; then
    echo "==> Installing errcheck..."
    go install github.com/kisielk/errcheck@latest
fi

err_files=$(errcheck -ignoretests \
                     -ignore 'github.com/hashicorp/terraform/helper/schema:Set' \
                     -ignore 'bytes:.*' \
                     -ignore 'io:Close|Write' \
                     $(go list ./...| grep -v /vendor/))

if [[ -n ${err_files} ]]; then
    echo 'Unchecked errors found in the following places:'
    echo "${err_files}"
    echo "Please handle returned errors. You can check directly with \`just errcheck\`"
    exit 1
fi

exit 0
