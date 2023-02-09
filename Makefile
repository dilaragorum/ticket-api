test:
	go test -v ./... -coverprofile=unit_coverage.out
	go tool cover -html=unit_coverage.out -o unit_coverage.html
	open unit_coverage.html

lint:
	golangci-lint run -c .golangci.yml -v

lint-fix:
	golangci-lint run -c .golangci.yml -v --fix

generate-mocks:
	mockgen -source internal/ticket/repository/repository.go -destination internal/ticket/mocks/repository.go -package mocks
	mockgen -source internal/ticket/service/service.go -destination internal/ticket/mocks/service.go -package mocks

docker-build:
	docker build -t api .

docker-run:
	docker run --rm -i -t -p 3000:3000 api