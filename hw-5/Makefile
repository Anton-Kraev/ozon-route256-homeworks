.PHONY: app
app:
	$(info running app...)
	docker-compose -f docker-compose.yml up --build -d postgres app

.PHONY: migrations
migrations:
	$(info running migrations...)
	./migration.sh

.PHONY: test
test:
	make unit-tests
	make integration-tests

.PHONY: unit-tests
unit-tests:
	$(info running unit tests...)
	go test ./internal/... ./pkg/...

.PHONY: integration-tests
integration-tests:
	$(info running integration tests...)
	go test ./tests/...

.PHONY: up-test-db
up-test-db:
	$(info running test db...)
	docker-compose -f docker-compose-test.yml up --build

.PHONY: up-test-migrations
up-test-migrations:
	$(info running migrations on test db...)
	./migration.sh test
