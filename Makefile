test:
	CI=true go test -v ./...
build:
	go build -o go-log-aggregator .
