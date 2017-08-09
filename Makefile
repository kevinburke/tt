.PHONY = test

SHELL = /bin/bash

BUMP_VERSION := $(GOPATH)/bin/bump_version
RELEASE := $(GOPATH)/bin/github-release

test:
	go vet ./...
	go test ./...

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

$(RELEASE):
	go get -u github.com/aktau/github-release

# Run "GITHUB_TOKEN=my-token make release version=0.x.y" to release a new version.
release: test | $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	$(BUMP_VERSION) --version=$(version) main.go
	git push origin --tags
	mkdir -p releases/$(version)
	# Change the binary names below to match your tool name
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/tt-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/tt-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/tt-windows-amd64 .
	# Change the Github username to match your username.
	# These commands are not idempotent, so ignore failures if an upload repeats
	$(RELEASE) release --user kevinburke --repo tt --tag $(version) || true
	$(RELEASE) upload --user kevinburke --repo tt --tag $(version) --name tt-linux-amd64 --file releases/$(version)/tt-linux-amd64 || true
	$(RELEASE) upload --user kevinburke --repo tt --tag $(version) --name tt-darwin-amd64 --file releases/$(version)/tt-darwin-amd64 || true
	$(RELEASE) upload --user kevinburke --repo tt --tag $(version) --name tt-windows-amd64 --file releases/$(version)/tt-windows-amd64 || true
