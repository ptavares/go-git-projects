###################
# RELEASE ENV
##################
RELEASE_SUPPORT=.release/make-release-support
RELEASE_FILE=.release/release
VERSION=$(shell . $(RELEASE_SUPPORT) ; getVersion)
TAG=$(shell . $(RELEASE_SUPPORT); getTag)
SHA=$(shell git show-ref -s $(TAG))

#####################
# COMMON VALUES
#####################
SHELL=/bin/bash
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")
MT = $(shell printf "  \033[36;1m▶\033[0m")
MT2 = $(shell printf "    \033[36;1m-\033[0m")

#####################
# APPLICATION VALUES
#####################
BINARY := git-projects
PKG := git-projects
LD_FLAGS := -ldflags="-X '$(PKG)/internal/version.Version=$(VERSION)'"
BUILD_FLAGS := $(LD_FLAGS) -v
ARCHIVE_DIR=archive
BUILD_DIR=build

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`


#####################
# TARGETS
#####################
default: .build version ; @ ## Default Task, build program with default values

build-windows: .build-windows version ; @ ## Build Go application for Windows OS

build-linux: .build-linux version ; @ ## Build Go application for Linux OS

build-darwin: .build-darwin version ; @ ## Build Go application for Mac Os

build-all: .build-all version ; @ ## Build Go application for all OS

archive: .do-archive version ; @ ## Build Go application archive for all OS

deps: .do-deps ; @ ## Download GO module dependency

fmt: .do-fmt ; @ ## Format GO code

lint: .do-lint ; @ ## Lint GO code

clean: .do-clean ; @ ## Clean GO app & directory

precommit: .precommit ; @ ## Execute some check with precommit hook

version: .do-version ; @ ## Get current version

check-status: .do-check-status ; @ ## Check current git status

check-release: .do-check-release ; @ ## Check release status

major-release: .do-major-release ; @ ## Do a major-release, ie : bumped first digit X+1.y.z

minor-release: .do-minor-release ; @ ## Do a minor-release, ie : bumped second digit x.Y+1.z

patch-release: .do-patch-release ; @ ## Do a patch-release, ie : bumped third digit x.y.Z+1

help: .do-help ; @ ## Show this help (Run make <target> V=1 to enable verbose)

# =====================
# ====   BUILD    =====
# =====================
.do-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/go-critic/go-critic/cmd/gocritic@latest
	go get github.com/google/go-github/v43
	go get github.com/google/go-github/v43/github
	go mod download
	go mod tidy


.build: .build-info .pre-build .do-build .post-build
.build-info: ; $(info $(M) Building… )

.pre-build: .do-fmt .do-lint

.do-build: ; $(info $(MT) Go Build)
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY) $(PKG)

.post-build:

.build-windows: ; $(info $(MT) Go Build for windows)
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/amd64/$(BINARY).exe $(PKG)
	env GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/arm64/$(BINARY).exe $(PKG)

.build-linux: ; $(info $(MT) Go Build for linux)
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/amd64/$(BINARY)-linux $(PKG)
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/arm64/$(BINARY)-linux $(PKG)

.build-darwin: ; $(info $(MT) Go Build for Mac)
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/amd64/$(BINARY)-darwin $(PKG)
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/arm64/$(BINARY)-darwin $(PKG)

.build-all: .build-windows .build-linux .build-darwin

.do-archive: .build-all
	$(info $(MT) Building archive...)
	cd $(BUILD_DIR)/amd64; tar -czf $(BINARY)-$(VERSION)-windows-amd64.tar.gz $(BINARY).exe
	cd $(BUILD_DIR)/arm64; tar -czf $(BINARY)-$(VERSION)-windows-arm64.tar.gz $(BINARY).exe
	cd $(BUILD_DIR)/amd64; tar -czf $(BINARY)-$(VERSION)-linux-amd64.tar.gz $(BINARY)-linux
	cd $(BUILD_DIR)/arm64; tar -czf $(BINARY)-$(VERSION)-linux-arm64.tar.gz $(BINARY)-linux
	cd $(BUILD_DIR)/amd64; tar -czf $(BINARY)-$(VERSION)-darwin-amd64.tar.gz $(BINARY)-darwin
	cd $(BUILD_DIR)/arm64; tar -czf $(BINARY)-$(VERSION)-darwin-arm64.tar.gz $(BINARY)-darwin
	[ -d $(ARCHIVE_DIR) ] || mkdir -p $(ARCHIVE_DIR)
	mv $(BUILD_DIR)/**/*.tar.gz $(ARCHIVE_DIR)/


# =====================
# ====  RELEASES  =====
# =====================

# ===> Major
.do-major-release: .major-release-info .tag-major-release .do-release version
.major-release-info: ; $(info $(M) Do major release...)
.tag-major-release: VERSION := $(shell . $(RELEASE_SUPPORT); nextMajorLevel)
.tag-major-release: .release .tag

# ===> Minor
.do-minor-release: .minor-release-info .tag-minor-release .do-release version
.minor-release-info: ; $(info $(M) Do minor release...)
.tag-minor-release: VERSION := $(shell . $(RELEASE_SUPPORT); nextMinorLevel)
.tag-minor-release: .release .tag

# ===> Path
.do-patch-release: .patch-release-info .tag-patch-release .do-release version
.patch-release-info: ; $(info $(M) Do minor release...)
.tag-patch-release: VERSION := $(shell . $(RELEASE_SUPPORT); nextPatchLevel)
.tag-patch-release: .release .tag


# ===> INIT RELEASE FILE
.release:
	@echo "release=0.0.0" > $(RELEASE_FILE)
	@echo "tag=0.0.0" >> $(RELEASE_FILE)
	@echo INFO: $(RELEASE_FILE) created
	@cat $(RELEASE_FILE)

# ===> DO RELEASE
.do-release: check-status check-release

# ===> Do TAG
.tag: TAG=$(shell . $(RELEASE_SUPPORT); getTag $(VERSION))
.tag: check-status
	@. $(RELEASE_SUPPORT) ; ! tagExists $(TAG) || (echo "ERROR: tag $(TAG) for version $(VERSION) already tagged in git" >&2 && exit 1) ;
	@. $(RELEASE_SUPPORT) ; setRelease $(VERSION)
	sed -i -e "s/Release_version-.*-blue/Release_version-$(VERSION)-blue/g" README.md
	sed -i -e "s/\"version\": \".*\"/\"version\": \"$(VERSION)\"/g" package.json
# sed -i -e "s/ref=.*\"/ref=$(VERSION)\"/g" examples/main.tf
	docker container run -it -v ${PWD}:/app --rm yvonnick/gitmoji-changelog:latest update $(VERSION)
	git add --all
	git commit -m ":bookmark: bumped to version $(VERSION)" ;
	git tag $(TAG) ;
	@ if [ -n "$(shell git remote -v)" ] ; then git push --tags ; else echo 'no remote to push tags to' ; fi
	git push

# ===> CHECK RELEASE
.do-check-release: ; $(info $(M) checking release...)
	@. $(RELEASE_SUPPORT) ; tagExists $(TAG) || (echo "ERROR: version not yet tagged in git. make [minor,major,patch]-release." >&2 && exit 1) ;
	@. $(RELEASE_SUPPORT) ; ! differsFromRelease $(TAG) || (echo "ERROR: current directory differs from tagged $(TAG). make [minor,major,patch]-release." ; exit 1)


# =======================
# ===    COMMONS    =====
# =======================

.precommit: ; $(info $(M) Checking precommit hook)
	pre-commit run -a

.do-fmt: ; $(info $(M) Formating code)
	go fmt

.do-lint: ; $(info $(M) Linting code)
	golangci-lint run

.do-clean:
	go clean
	rm -rf $(BUILD_DIR)/ $(ARCHIVE_DIR)

# ===> Get current version
.do-version: ; $(info $(M) current version)
	$(info $(MT) $(VERSION))

# ===> Check current repository status
.do-check-status:  ; $(info $(M) checking git status...)
	@. $(RELEASE_SUPPORT) ; ! hasChanges || (echo "ERROR: there are still outstanding changes" >&2 && exit 1) ;

# ===> Help
.do-help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# ===========================================================================================

BUILD_TARGETS := pre-build do-build post-build build lint fmt clean deps
REALEASE_TARGETS := check-release major-release  minor-release patch-release
INFO_TARGETS := version .do-version check-status .do-check-status precommit .precommit
HELP_TARGETS :=  help .do-help

.PHONY: $(BUILD_TARGETS) $(RELEASE_TARGETS) $(INFO_TARGETS) $(HELP_TARGETS)
