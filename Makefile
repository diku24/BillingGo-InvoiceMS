build:
	@go build -o bin/InvoiceMS

run: build
	@./bin/InvoiceMS

test:
	@go test ./... -coverprofile=coverage

testHTMLCoverage: test
	@go tool cover -html=coverage

testFuncCoverage: test
	@go tool cover -func=coverage

mody:
	@go mod tidy

resetMody:
	@go mod tidy -v

# docker compose watch
up:
	docker compose up

down:
	docker compose down

prune:
	docker image prune -f

log:
	docker compose logs -f -t