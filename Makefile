TEST?=$$(go list ./... | grep -v vendor)
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOGEN_FILES?=$$(go list ./... | grep -v vendor)
BIN_NAME?=slack-bot-template
GO_FILES?=$(shell find . -name '*.go')

BUILD_LDFLAGS = " -s -w "

.PHONY: default
default: test build

.PHONY: run
run:
	go run $(CURDIR)/main.go $(ARGS)

.PHONY: clean
clean:
	rm -Rf bin/*

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep; \
	go get -u github.com/golang/lint/golint


.PHONY: build
build: bin/slack-bot-template

bin/slack-bot-template: $(GO_FILES)
	GOOS="`go env GOOS`" GOARCH="`go env GOARCH`" CGO_ENABLED=0 \
	go build -ldflags $(BUILD_LDFLAGS) -o bin/$(BIN_NAME)

.PHONY: test
test: vet
	go test $(TEST) $(TESTARGS) -v -timeout=30m -parallel=4 ;

.PHONY: vet
vet: golint
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: golint
golint: fmt
	for pkg in $$(go list ./... | grep -v /vendor/ ) ; do \
        test -z "$$(golint $$pkg | grep -v '_gen.go' | grep -v '_string.go' | grep -v 'should have comment' | tee /dev/stderr)" || RES=1; \
    done ;exit $$RES

.PHONY: fmt
fmt:
	gofmt -s -l -w $(GOFMT_FILES)
