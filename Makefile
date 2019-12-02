GOOS = $(word 2, $(subst _, ,$(basename $(notdir $(@)))))
GOARCH = $(word 3, $(subst _, ,$(basename $(notdir $(@)))))
LINTER ?= golangci-lint
LINTER_ARGS ?= run

dist:
	mkdir -p $(@)

dist/%: | dist
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(@)

.DEFAULT_GOAL := all
.PHONY: all
all:  dist/pickerbot_darwin_amd64 \
	dist/pickerbot_linux_amd64 \
	dist/pickerbot_linux_arm64 \
	dist/pickerbot_windows_amd64.exe

.PHONY: clean
clean:
	rm -rf dist

.PHONY: lint
lint:
	$(LINTER) $(LINTER_ARGS)

.PHONY: test
test:
	go test $(TEST_ARGS)