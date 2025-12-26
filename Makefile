build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

docker:
	@echo "building docker image"
	@docker build -t api .
	@echo "running API inside Docker container"
	@docker run -p 3000:3000 api

test:
	@MONGO_DB_URL_TEST=mongodb://localhost:27017 go test -v ./...