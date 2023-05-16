run:
	go run .

run-compose:
	COMPOSE_PROJECT_NAME=garnbarn docker-compose up -d

mockgen:
	mockgen -destination=./test/mock_repository/mock_tag.go -source=./repository/tag.go -package=mock_repository
	mockgen -destination=./test/mock_service/mock_service -source=./service/tag.go -package=mock_service