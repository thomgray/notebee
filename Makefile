run:
	notebeedevmode=true go run .

build:
	go build -o target/notebee .code

install:
	go install

test:
	go test ./...