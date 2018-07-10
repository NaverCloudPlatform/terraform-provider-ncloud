TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
EXEC_FILE=terraform-provider-ncloud_v$(VERSION)
PKG_NAME=ncloud

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

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

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

build_all_platforms:
	GOOS=linux GOARCH=amd64 go build -o $(EXEC_FILE) && zip terraform-provider-ncloud_linux_amd64_$(VERSION).zip $(EXEC_FILE) && rm $(EXEC_FILE)
	GOOS=linux GOARCH=386 go build -o $(EXEC_FILE) && zip terraform-provider-ncloud_linux_386_$(VERSION).zip $(EXEC_FILE) && rm $(EXEC_FILE)
	GOOS=linux GOARCH=arm go build -o $(EXEC_FILE) && zip terraform-provider-ncloud_linux_arm_$(VERSION).zip $(EXEC_FILE) && rm $(EXEC_FILE)
	GOOS=darwin GOARCH=386 go build -o $(EXEC_FILE) && zip terraform-provider-ncloud_darwin_386_$(VERSION).zip $(EXEC_FILE) && rm $(EXEC_FILE)
	GOOS=darwin GOARCH=amd64 go build -o $(EXEC_FILE) && zip terraform-provider-ncloud_darwin_amd64_$(VERSION).zip $(EXEC_FILE) && rm $(EXEC_FILE)
	GOOS=windows GOARCH=amd64 go build -o $(EXEC_FILE).exe && zip terraform-provider-ncloud_windows_amd64_$(VERSION).zip $(EXEC_FILE).exe && rm $(EXEC_FILE).exe
	GOOS=windows GOARCH=386 go build -o $(EXEC_FILE).exe && zip terraform-provider-ncloud_windows_386_$(VERSION).zip $(EXEC_FILE).exe && rm $(EXEC_FILE).exe

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test

