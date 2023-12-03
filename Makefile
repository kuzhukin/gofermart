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
	cmd/gophermart/gophermart -a localhost:8080 -d fakeDbUri -r fakeAccuralAddr

generate:
	go generate ./...