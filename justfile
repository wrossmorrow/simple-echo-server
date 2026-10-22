
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
    docker build -t echo:{{tag}} .
