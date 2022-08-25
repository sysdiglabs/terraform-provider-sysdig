SWEEP?=us-east-1,us-west-2
TEST?=./...
PKG_NAME=sysdig
WEBSITE_REPO=github.com/hashicorp/terraform-website
VERSION=$(shell [ ! -z `git tag -l --contains HEAD` ] && git tag -l --contains HEAD || git rev-parse --short HEAD)
GOPATH=$(shell go env GOPATH)

default: build

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	CGO_ENABLED=1 TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -race

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
