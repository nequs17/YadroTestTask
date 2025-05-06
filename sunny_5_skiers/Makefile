BINARY_NAME=skiers

BIN_DIR=bin

BINARY_PATH=$(BIN_DIR)/$(BINARY_NAME)

MAIN_PACKAGE=cmd/main.go

BUILD_FLAGS=-o $(BINARY_PATH)

LINT_CMD=golangci-lint run

LINT_CONFIG=.golangci.yml

GO=go

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(BUILD_FLAGS) $(MAIN_PACKAGE)

run: build
	./$(BINARY_PATH)

clean:
	@rm -rf $(BIN_DIR)
	@rm -f output.log
	@rm -f results.txt

lint:
	$(LINT_CMD) --config=$(LINT_CONFIG) ./...

lint-fix:
	$(LINT_CMD) --config=$(LINT_CONFIG) --fix ./...

test:
	$(GO) test -v -cover ./...
 
.PHONY: test, lint, lint-fix, clean, run, build, all