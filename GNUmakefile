SWEEP?=us-east-1,us-west-2
TEST?=./...
TEST_SUITE?=tf_acc_sysdig_monitor,tf_acc_sysdig_secure
PKG_NAME=sysdig
WEBSITE_REPO=github.com/hashicorp/terraform-website
VERSION=$(shell [ ! -z `git tag -l --contains HEAD` ] && git tag -l --contains HEAD || git rev-parse --short HEAD)
GOPATH=$(shell go env GOPATH)

TERRAFORM_PLUGIN_ROOT_DIR=$(HOME)/.terraform.d/plugins
TERRAFORM_PROVIDER_REFERENCE_NAME=local
TERRAFORM_PROVIDER_NAME=sysdiglabs/$(PKG_NAME)
TERRAFORM_PROVIDER_DEV_VERSION=1.0.0
TERRAFORM_PLATFORM=$(shell terraform version -json | jq -r .platform)
TERRAFORM_SYSDIG_PLUGIN_DIR=$(TERRAFORM_PLUGIN_ROOT_DIR)/$(TERRAFORM_PROVIDER_REFERENCE_NAME)/$(TERRAFORM_PROVIDER_NAME)/$(TERRAFORM_PROVIDER_DEV_VERSION)/$(TERRAFORM_PLATFORM)

install-tools:
	go install golang.org/x/tools/cmd/stringer@latest

default: build

build: fmtcheck
	go install

install: fmtcheck
	go build -o terraform-provider-sysdig
	mkdir -p $(TERRAFORM_SYSDIG_PLUGIN_DIR)
	cp terraform-provider-sysdig $(TERRAFORM_SYSDIG_PLUGIN_DIR)/terraform-provider-sysdig

uninstall:
	rm -rf $(TERRAFORM_SYSDIG_PLUGIN_DIR)


sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -tags=unit -timeout=30s -parallel=4

testacc: fmtcheck
	CGO_ENABLED=1 TF_ACC=1 go test $(TEST) -v $(TESTARGS) -tags=$(TEST_SUITE) -timeout 120m -race -parallel=1

junit-report: fmtcheck
	echo "Current directory $PWD"
	@go install github.com/jstemmer/go-junit-report/v2@latest
	CGO_ENABLED=1 TF_ACC=1 TF_LOG=DEBUG go test $(TEST) -v $(TESTARGS) -tags=$(TEST_SUITE) -timeout 120m -race -parallel=1 2>&1 | tee output.txt
	! grep -q "\[build failed\]" output.txt
	go-junit-report -in output.txt -out junit-report.xml

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	go install mvdan.cc/gofumpt@latest
	go fmt ./...
	gofumpt -w ./

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	golangci-lint run --timeout 1h ./...

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

release: fmtcheck
	for kernel in linux windows darwin; do \
		for dist in $$(go tool dist list | grep $$kernel); do  \
			GOOS=$$kernel; \
			GOARCH=$$(echo $$dist | cut -d/ -f2); \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -o terraform-provider-sysdig_$(VERSION); \
			tar -czf terraform-provider-sysdig-$$GOOS-$$GOARCH.tar.gz terraform-provider-sysdig_$(VERSION) --remove-files; \
		done \
	done

.PHONY: website
website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	(cd $(GOPATH)/src/$(WEBSITE_REPO); \
	  ln -s $(shell pwd) ext/providers/sysdig; \
	  ln -s ../../../ext/providers/sysdig/website/sysdig.erb content/source/layouts/sysdig.erb; \
	  ln -s ../../../../ext/providers/sysdig/website/docs content/source/docs/providers/sysdig; \
	)
endif
	$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: website-test
website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	(cd $(GOPATH)/src/$(WEBSITE_REPO); \
	  ln -s $(shell pwd) ext/providers/sysdig; \
	  ln -s ../../../ext/providers/sysdig/website/sysdig.erb content/source/layouts/sysdig.erb; \
	  ln -s ../../../../ext/providers/sysdig/website/docs content/source/docs/providers/sysdig; \
	)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)
