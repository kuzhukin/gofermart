ALL_TARGETS = gophermart

define build
	go build -o ./cmd/$(1)/ -v ./cmd/$(1)/ 
endef

all: build

.PHONY: build
build: $(patsubst %, build-%, $(ALL_TARGETS))


.PHONY: build-%
build-%:
	@echo === Building $*
	$(call build,$*)

.PHONY: test
test:
	@echo === Tests
	go test -count 1 -v -cover ./...

define clean
	rm ./cmd/$(1)/$(1)
endef

.PHONY: clean
clean: $(patsubst %, clean-%, $(ALL_TARGETS))

clean-%:
	@echo === Cleaning $*
	$(call clean,$*)

# Linter constants
LINTER := golangci-lint 

.PHONY: lint
lint:
	@echo === Lint
	$(LINTER) --version
	$(LINTER) cache clean && $(LINTER) run

run:
	docker stop gophermartpostgres && \
	docker rm gophermartpostgres && \
	docker run -d --name gophermartpostgres -p 5431:5432 -e POSTGRES_DB=gophermart -e POSTGRES_USER=gophermart -e POSTGRES_PASSWORD=12345 postgres:12-alpine && \
	sleep 5 && \
	cmd/gophermart/gophermart -a localhost:8080 -d "postgres://gophermart:12345@127.0.0.1:5431/gophermart" -r "localhost:33555"

run_accrual:
	cmd/accrual/accrual_linux_amd64 -a "localhost:33555"

generate:
	go generate ./...