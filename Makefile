SHELL = /bin/bash

BUMP_VERSION := $(shell command -v bump_version)
RELEASE := $(shell command -v github-release)

test:
	go vet ./...
	go test ./...

# Run "GITHUB_TOKEN=my-token make release version=0.x.y" to release a new version.
release: test
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
ifndef BUMP_VERSION
	go get -u github.com/Shyp/bump_version
endif
	bump_version --version=$(version) main.go
	git push origin --tags
	mkdir -p releases/$(version)
	# Change the binary names below to match your tool name
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/tt-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/tt-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/tt-windows-amd64 .
ifndef RELEASE
	go get -u github.com/aktau/github-release
endif
	# Change the Github username to match your username.
	# These commands are not idempotent, so ignore failures if an upload repeats
	github-release release --user kevinburke --repo tt --tag $(version) || true
	github-release upload --user kevinburke --repo tt --tag $(version) --name tt-linux-amd64 --file releases/$(version)/tt-linux-amd64 || true
	github-release upload --user kevinburke --repo tt --tag $(version) --name tt-darwin-amd64 --file releases/$(version)/tt-darwin-amd64 || true
	github-release upload --user kevinburke --repo tt --tag $(version) --name tt-windows-amd64 --file releases/$(version)/tt-windows-amd64 || true
