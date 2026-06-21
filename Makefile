run:
	docker compose up

unit:
	go test ./internal/...

integration:
	docker compose -f ./tests/docker/compose.yaml up