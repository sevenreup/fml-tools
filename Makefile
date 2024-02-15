build:
	@go build -C ./src -o ../bin/fml

run: build
	@./bin/fml

test:
	@go test -v ./...