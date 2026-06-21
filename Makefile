run:
	docker compose up

unit:
	go test ./internal/... -coverprofile=coverage_unit.out -coverpkg=github.com/Sugyk/auth_service/...

integration:
	docker compose -f ./tests/docker/compose.yaml up -d
	go test ./tests/... -coverprofile=coverage_integrational.out -coverpkg=github.com/Sugyk/auth_service/...
	docker compose -f ./tests/docker/compose.yaml down

cover:
	go tool cover -func=$(f)
	go tool cover -html=$(f)
