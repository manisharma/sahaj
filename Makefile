coverage:
	go test -timeout 30s -coverprofile=go-code-cover ./...

test:
	go test -v -timeout 30s ./...

run:
	go run main.go