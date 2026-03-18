include .env
export

.PHONY: build_server start build_test_migrates

build_server:
	go build -o ./build/server -v ./cmd/server

start: build_server
	$(HOME)/go/bin/migrate -path ./migrations -database "$(DB_URL)" drop
	$(HOME)/go/bin/migrate -path ./migrations -database "$(DB_URL)" up
	./build/server

build_test_migrates:
	$(HOME)/go/bin/migrate -path ./migrations -database "$(TEST_DB_URL)" drop
	$(HOME)/go/bin/migrate -path ./migrations -database "$(TEST_DB_URL)" up