# Targets not related to individual files
.PHONY: all test cover_func cover_html out clean vet loc fmt test_v bench

# Build constants
BUILD_OUT_DIR = bin
TEST_COVERAGE_PROFILE = coverage.out

all: out fmt vet test_v loc

out:
	mkdir -p $(BUILD_OUT_DIR)

fmt:
	go fmt ./...

vet:
	go vet ./...

test: out
	go test -race ./... -coverprofile=$(BUILD_OUT_DIR)/$(TEST_COVERAGE_PROFILE)

test_v: out
	go test -race -v ./... -coverprofile=$(BUILD_OUT_DIR)/$(TEST_COVERAGE_PROFILE)

bench:
	go test -bench=.

loc:
	find . -type f -not -path "./vendor/*" -name "*.go" | xargs wc -l

cover_func: test
	go tool cover -func=$(BUILD_OUT_DIR)/$(TEST_COVERAGE_PROFILE)

cover_html: test
	go tool cover -html=$(BUILD_OUT_DIR)/$(TEST_COVERAGE_PROFILE)

clean:
	rm -rf $(BUILD_OUT_DIR)
