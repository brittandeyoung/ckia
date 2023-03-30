GO_VER              ?= go

default: build

build:
	$(GO_VER) install

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s .

test:
	@echo "==> Running Unit Tests..."
	go test ./... -race -covermode=atomic -coverprofile=coverage.out


pr: test build

changelog-release:
	pwsh -noprofile -command  'Update-Changelog -ReleaseVersion $(RELEASE_VERSION) -LinkMode Automatic -LinkPattern @{ FirstRelease = "https://github.com/brittandeyoung/ckia/tree/v\{CUR\}"; NormalRelease = "https://github.com/brittandeyoung/ckia/compare/v\{PREV\}..v\{CUR\}"; Unreleased = "https://github.com/brittandeyoung/ckia/compare/v\{CUR\}..HEAD"}' 

changelog-add:
	pwsh -noprofile -command  'Add-ChangelogData -Type "$(TYPE)" -Data "$(DATA)"'

.PHONY: \
	build \
	changelog-release \
	changelog-add \
	fmt \
	pr \
	test \
