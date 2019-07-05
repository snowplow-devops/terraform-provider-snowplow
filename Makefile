.PHONY: all format lint test release release-dry dep clean

# -----------------------------------------------------------------------------
#  CONSTANTS
# -----------------------------------------------------------------------------

version = `cat VERSION`

src_dir := terraform-provider-snowplow

build_dir = build

coverage_dir  = $(build_dir)/coverage
coverage_out  = $(coverage_dir)/coverage.out
coverage_html = $(coverage_dir)/coverage.html

output_dir    = $(build_dir)/output

linux_dir     = $(output_dir)/linux
darwin_dir    = $(output_dir)/darwin
windows_dir   = $(output_dir)/windows

bin_name      = terraform-provider-snowplow_v$(version)
bin_linux     = $(linux_dir)/$(bin_name)
bin_darwin    = $(darwin_dir)/$(bin_name)
bin_windows   = $(windows_dir)/$(bin_name)

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all: dep
	go get -u github.com/mitchellh/gox/...
	gox -osarch=linux/amd64 -output=$(bin_linux) ./$(src_dir)
	gox -osarch=darwin/amd64 -output=$(bin_darwin) ./$(src_dir)
	gox -osarch=windows/amd64 -output=$(bin_windows) ./$(src_dir)

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	go fmt ./$(src_dir)
	gofmt -s -w ./$(src_dir)

lint:
	go get -u golang.org/x/lint/golint
	golint ./$(src_dir)

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test: dep
	mkdir -p $(coverage_dir)
	go get -u golang.org/x/tools/cmd/cover/...
	go test ./$(src_dir) -tags test -v -covermode=count -coverprofile=$(coverage_out)
	go tool cover -html=$(coverage_out) -o $(coverage_html)

# -----------------------------------------------------------------------------
#  RELEASE
# -----------------------------------------------------------------------------

release:
	release-manager --config .release.yml --check-version --make-artifact --make-version --upload-artifact

release-dry:
	release-manager --config .release.yml --check-version --make-artifact

# -----------------------------------------------------------------------------
#  DEPENDENCIES
# -----------------------------------------------------------------------------

dep:
	dep ensure

# -----------------------------------------------------------------------------
#  CLEANUP
# -----------------------------------------------------------------------------

clean:
	rm -rf $(build_dir)
