.PHONY: run proto swagger unit integration cover

run:
	docker image prune -f
	docker compose build
	docker compose up
	docker compose down

proto:
	protoc \
		--go_out=. --go_opt=module=github.com/Sugyk/auth_service \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Sugyk/auth_service \
		proto/auth.proto

swagger:
	swag init -g cmd/auth_service/main.go -o docs --parseDependency --parseInternal

unit:
	go test ./internal/... -coverprofile=coverage_unit.out -coverpkg=github.com/Sugyk/auth_service/...

integration:
	docker compose -f ./tests/docker/compose.yaml up -d
	go test ./tests/... -coverprofile=coverage_integration.out -coverpkg=github.com/Sugyk/auth_service/...
	docker compose -f ./tests/docker/compose.yaml down

cover:
	go tool cover -func=$(f)
	go tool cover -html=$(f)
