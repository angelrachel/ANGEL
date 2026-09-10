.PHONY: fmt test vet build check inventory static-check safe-contracts frontend-build gateway-test planner-test lab-check lab-smoke app-check validate

fmt:
	files="$$(gofmt -l .)"; test -z "$$files" || gofmt -w $$files

test:
	go test ./src/...

vet:
	go vet ./src/...

build:
	go build ./src/...

inventory:
	python3 scripts/generate_module_inventory.py

static-check:
	python3 scripts/static_success_checker.py

safe-contracts:
	python3 scripts/validate_safe_contracts.py

frontend-build:
	npm --prefix apps/frontend-angular run build

gateway-test:
	dotnet test apps/gateway-dotnet.tests/AngelGateway.Tests.csproj

planner-test:
	.venv-langgraph/bin/python -m unittest discover -s apps/orchestrator-langgraph -p 'test_*.py'

lab-check:
	cd lab/services/fixture-http && go test ./...

lab-smoke:
	bash lab/acceptance/p0_fixture_http.sh

app-check: frontend-build gateway-test planner-test lab-check lab-smoke

validate: check inventory static-check safe-contracts
	bash -n deploy/scripts/deploy.sh deploy/scripts/destroy.sh deploy/scripts/validate.sh
	python3 scripts/repo_integrity_check.py

check: fmt test vet build
