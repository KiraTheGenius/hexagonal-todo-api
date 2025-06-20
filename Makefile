.PHONY: run test benchmark clean setup-localstack

run:
	docker-compose up --build

test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem ./tests/...

clean:
	docker-compose down -v
	docker system prune -f

setup-localstack:
	@echo "Setting up LocalStack S3 bucket..."
	@docker-compose exec localstack awslocal s3 mb s3://todo-files

dev:
	docker-compose up mysql redis localstack -d
	sleep 5
	make setup-localstack
	go run cmd/server/main.go

mod:
	go mod tidy
	go mod download

lint:
	golangci-lint run