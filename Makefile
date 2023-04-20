BINARY_NAME=cb-tracker
BINARY_FILE := ./bin/$(BINARY_NAME)
GOTEST_DIR := test-results
GOTEST_FLAGS := -cover -race -v -count=1 -timeout 60s

VERSION ?= local
SCM_COMMIT ?= `git rev-parse HEAD`

.PHONY: build test run
build:
	@echo ">> building application"
	go build -trimpath -ldflags \
	"-X main.Version=$(VERSION) \
	-X main.SCMCommit=$(SCM_COMMIT)" \
	-o $(BINARY_FILE) \
	./cmd/...

run:
	go run ./cmd/cp-repl-ping

test:
	@echo ">> running all tests"
	go test $(GOTEST_FLAGS) ./...

test-deps:
	@which gotestsum > /dev/null || \
		go install gotest.tools/gotestsum@latest

test-with-junit-report: test-deps
	gotestsum --junitfile $(GOTEST_DIR)/tests.xml -- $(GOTEST_FLAGS) ./...

lint-deps:
	@which golangci-lint > /dev/null || \
		(curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.1)

lint: lint-deps
	golangci-lint run

clean:
	go clean
	rm -rf $(GOTEST_DIR)
	rm $(BINARY_FILE)

mocks-deps:
	@which mockgen > /dev/null || \
		(go get github.com/golang/mock/mockgen)

mocks: mocks-deps
	@echo "==> update mocks"
	@echo ""
	mockgen -source pkg/kv/interfaces/kv.go -package=mocks > pkg/kv/interfaces/mocks/kv.go
	mockgen -source pkg/pinger/interfaces/pinger.go -package=mocks > pkg/pinger/interfaces/mocks/pinger.go
