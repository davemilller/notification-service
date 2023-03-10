.PHONY: lint
lint:
	${INFO} "Linting..."
	staticcheck ./...
	go vet ./...
	
.PHONY: test
test:
	${INFO} "Running unit tests..."
	mkdir -p target
	go test -coverprofile target/coverage.out -race -count=10 -short ./...
	go tool cover -html=target/coverage.out -o target/coverage.html
	${INFO} "Unit testing complete"

.PHONY: test-integration
test-integration:
	${INFO} "Running unit and integration tests..."
	mkdir -p target
	go test -coverprofile target/coverage.out -race -count=1 ./...
	go tool cover -html=target/coverage.out -o target/coverage.html
	${INFO} "Integration testing complete"

.PHONY: lint-test
lint-test: lint
	test

YELLOW := "\e[1;33m"
GREEN := "\e[1;32m"
NC := "\e[0m"

INFO := @bash -c '\
	printf $(GREEN); \
	echo "=> $$1"; \
	printf $(NC)' VALUE
