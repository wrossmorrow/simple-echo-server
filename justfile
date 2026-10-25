
image_name := "ghcr.io/wrossmorrow/simple-echo-server"

# list targets and exit
default:
    just --list

# run the server
run:
    go run main.go

# build locally
build:
    go build -o bin/echo

# run tests
test:
    go test ./...

# containerize
containerize tag="local":
    docker build -t {{image_name}}:{{tag}} .
