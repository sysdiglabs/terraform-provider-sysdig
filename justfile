sweep := env_var_or_default("SWEEP", "us-east-1,us-west-2")
sweep_args := env_var_or_default("SWEEPARGS", "")
test := env_var_or_default("TEST", "./...")
test_args := env_var_or_default("TESTARGS", "")
test_suite := env_var_or_default("TEST_SUITE", "tf_acc_sysdig_monitor,tf_acc_sysdig_secure")
pkg_name := "sysdig"
website_repo := "github.com/hashicorp/terraform-website"

# List available recipes. Runs when `just` is called with no recipe.
[private]
default:
    @just --list

# Install dev tool dependencies
[group('setup')]
install-tools:
    go install golang.org/x/tools/cmd/stringer@latest

# Update nix flake inputs, Go module dependencies, pinned GitHub Actions, and pre-commit hook versions
[group('update')]
update:
    nix flake update
    nix develop --command go get -u -t -v ./...
    nix develop --command go mod tidy
    nix develop --command pinact run -u
    nix develop --command prek autoupdate

# Build and install the provider binary to $GOPATH/bin
[group('build')]
build: fmtcheck
    go install

# Build and install the provider to the local Terraform plugins directory
[group('install')]
install: fmtcheck
    #!/usr/bin/env bash
    set -euo pipefail
    go build -o terraform-provider-sysdig
    platform=$(terraform version -json | jq -r .platform)
    plugin_dir="$HOME/.terraform.d/plugins/local/sysdiglabs/{{ pkg_name }}/1.0.0/$platform"
    mkdir -p "$plugin_dir"
    cp terraform-provider-sysdig "$plugin_dir/terraform-provider-sysdig"

# Remove the provider from the local Terraform plugins directory
[group('install')]
uninstall:
    #!/usr/bin/env bash
    set -euo pipefail
    platform=$(terraform version -json | jq -r .platform)
    rm -rf "$HOME/.terraform.d/plugins/local/sysdiglabs/{{ pkg_name }}/1.0.0/$platform"

# WARNING: destroys infrastructure. Use only in development accounts.
[group('test')]
[confirm("This will destroy infrastructure. Use only in development accounts. Continue?")]
sweep:
    go test {{ test }} -v -sweep={{ sweep }} {{ sweep_args }}

# Run unit tests
[group('test')]
test: fmtcheck
    go test {{ test }} -tags=unit -timeout=30s -parallel=4

# Run acceptance tests (requires credentials)
[group('test')]
testacc: fmtcheck
    CGO_ENABLED=1 TF_ACC=1 go test {{ test }} -v {{ test_args }} -tags={{ test_suite }} -timeout 120m -race -parallel=1

# Run acceptance tests with debug logging and produce a JUnit report
[group('test')]
junit-report: fmtcheck
    #!/usr/bin/env bash
    set -eu
    go install github.com/jstemmer/go-junit-report/v2@latest
    CGO_ENABLED=1 TF_ACC=1 TF_LOG=DEBUG go test {{ test }} -v {{ test_args }} -tags={{ test_suite }} -timeout 120m -race -parallel=1 2>&1 | tee output.txt
    ! grep -q "\[build failed\]" output.txt
    go-junit-report -in output.txt -out junit-report.xml

# Run go vet across all non-vendor packages
[group('quality')]
vet:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "go vet ."
    if ! go vet $(go list ./... | grep -v vendor/); then
        echo ""
        echo "Vet found suspicious constructs. Please check the reported constructs"
        echo "and fix them if necessary before submitting the code for review."
        exit 1
    fi

# Format code with go fmt and gofumpt
[group('quality')]
fmt:
    go install mvdan.cc/gofumpt@latest
    go fmt ./...
    gofumpt -w ./

# Verify code formatting
[group('quality')]
fmtcheck:
    @./scripts/gofmtcheck.sh

# Run golangci-lint across all build tags
[group('quality')]
lint:
    golangci-lint run --build-tags "{{ test_suite }}" --timeout 1h ./...

# Build the provider, generate its schema, and validate docs against it
[group('docs')]
provider-docs:
    #!/usr/bin/env bash
    set -euo pipefail
    goos=$(go env GOOS)
    goarch=$(go env GOARCH)
    go build -o "terraform-plugin-dir/registry.terraform.io/sysdiglabs/sysdig/99.99.99/${goos}_${goarch}/terraform-provider-sysdig" .
    printf 'terraform {\n  required_providers {\n    sysdig = { source = "sysdiglabs/sysdig" }\n  }\n}\n' > main.tf
    terraform init -plugin-dir terraform-plugin-dir
    mkdir -p terraform-providers-schema
    terraform providers schema -json > terraform-providers-schema/schema.json
    go install github.com/bflad/tfproviderdocs@latest
    tfproviderdocs check \
        -allowed-resource-subcategories-file website/allowed-subcategories.txt \
        -enable-contents-check \
        -provider-source registry.terraform.io/sysdiglabs/sysdig \
        -providers-schema-json terraform-providers-schema/schema.json \
        -require-resource-subcategory
    rm -f main.tf
    rm -rf terraform-plugin-dir terraform-providers-schema .terraform .terraform.lock.hcl

# Check for unchecked errors
[group('quality')]
errcheck:
    @./scripts/errcheck.sh

# Show status of vendored dependencies
[group('utils')]
vendor-status:
    @govendor status

# Compile tests without running them; set TEST to a specific package first
[group('test')]
test-compile:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{ test }}" = "./..." ]; then
        echo "ERROR: Set TEST to a specific package. For example,"
        echo "  TEST=./sysdig just test-compile"
        exit 1
    fi
    go test -c {{ test }} {{ test_args }}

# Build and package release binaries for all supported platforms
[group('release')]
release: fmtcheck
    #!/usr/bin/env bash
    set -euo pipefail
    version=$([ -n "$(git tag -l --contains HEAD)" ] && git tag -l --contains HEAD || git rev-parse --short HEAD)
    for kernel in linux windows darwin; do
        for dist in $(go tool dist list | grep "$kernel"); do
            goarch=$(echo "$dist" | cut -d/ -f2)
            GOOS="$kernel" GOARCH="$goarch" go build -o "terraform-provider-sysdig_${version}"
            tar -czf "terraform-provider-sysdig-${kernel}-${goarch}.tar.gz" "terraform-provider-sysdig_${version}" --remove-files
        done
    done

# Serve provider docs locally via the (legacy) HashiCorp terraform-website repo
[group('docs')]
website:
    #!/usr/bin/env bash
    set -euo pipefail
    gopath=$(go env GOPATH)
    website_dir="$gopath/src/{{ website_repo }}"
    if [ ! -d "$website_dir" ]; then
        echo "{{ website_repo }} not found in your GOPATH (necessary for layouts and assets), getting..."
        git clone "https://{{ website_repo }}" "$website_dir"
        ln -s "$(pwd)" "$website_dir/ext/providers/sysdig"
        ln -s ../../../ext/providers/sysdig/website/sysdig.erb "$website_dir/content/source/layouts/sysdig.erb"
        ln -s ../../../../ext/providers/sysdig/website/docs "$website_dir/content/source/docs/providers/sysdig"
    fi
    make -C "$website_dir" website-provider PROVIDER_PATH="$(pwd)" PROVIDER_NAME={{ pkg_name }}

# Run the provider docs website's own test suite via the legacy terraform-website repo
[group('docs')]
website-test:
    #!/usr/bin/env bash
    set -euo pipefail
    gopath=$(go env GOPATH)
    website_dir="$gopath/src/{{ website_repo }}"
    if [ ! -d "$website_dir" ]; then
        echo "{{ website_repo }} not found in your GOPATH (necessary for layouts and assets), getting..."
        git clone "https://{{ website_repo }}" "$website_dir"
        ln -s "$(pwd)" "$website_dir/ext/providers/sysdig"
        ln -s ../../../ext/providers/sysdig/website/sysdig.erb "$website_dir/content/source/layouts/sysdig.erb"
        ln -s ../../../../ext/providers/sysdig/website/docs "$website_dir/content/source/docs/providers/sysdig"
    fi
    make -C "$website_dir" website-provider-test PROVIDER_PATH="$(pwd)" PROVIDER_NAME={{ pkg_name }}
