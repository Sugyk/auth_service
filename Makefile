run:
	docker image prune -f
	docker compose build
	docker compose up
	docker compose down

unit:
	go test ./internal/... -coverprofile=coverage_unit.out -coverpkg=github.com/Sugyk/auth_service/...

integration:
	docker compose -f ./tests/docker/compose.yaml up -d
	go test ./tests/... -coverprofile=coverage_integration.out -coverpkg=github.com/Sugyk/auth_service/...
	docker compose -f ./tests/docker/compose.yaml down

cover:
	go tool cover -func=$(f)
	go tool cover -html=$(f)
