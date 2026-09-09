.PHONY: fmt test vet build check inventory static-check validate

fmt:
	files="$$(gofmt -l .)"; test -z "$$files" || gofmt -w $$files

test:
	go test ./...

vet:
	go vet ./...

build:
	go build ./...

inventory:
	python3 scripts/generate_module_inventory.py

static-check:
	python3 scripts/static_success_checker.py

validate: check inventory static-check
	bash -n deploy/scripts/deploy.sh deploy/scripts/destroy.sh deploy/scripts/validate.sh
	python3 scripts/repo_integrity_check.py

check: fmt test vet build
