.PHONY: clean vet test

APP_NAME = raspall
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build ./tmp

vet:
	go vet ./...

test:
	gotest -v ./...

build: test vet
	go build -o $(BUILD_DIR)/$(APP_NAME) -v ./cmd/$(APP_NAME)

run: build
	$(BUILD_DIR)/$(APP_NAME)

dev: test vet
	air

