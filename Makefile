SHELL=/usr/bin/env bash

# Project specific properties.
application_name        = rosenbridge
application_binary_name = main

# Container specific properties.
application_image_name     = rosenbridge
application_container_name = rosenbridge-1

# For ProtoBuf code generation.
proto_path=pkg/proto/*.proto

# Builds the project.
build:
	@echo "+$@"
	@go build -o bin/$(application_binary_name) cmd/$(application_name)/main.go

# Runs the project after linting and building it anew.
run: tidy lint build
	@echo "+$@"
	@echo "########### Running the application binary ############"
	@bin/$(application_binary_name)

# Runs the project.
run-quick:
	@echo "+$@"
	@echo "########### Running the application binary ############"
	@bin/$(application_binary_name)

# Tests the whole project.
test:
	@echo "+$@"
	@CGO_ENABLED=1 go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Runs the "go mod tidy" command.
tidy:
	@echo "+$@"
	@go mod tidy

# Runs golang-ci-lint over the project.
lint:
	@echo "+$@"
	@golangci-lint run ./...

# Builds the docker image for the project.
image:
	@echo "+$@"
	@docker build --network host --file Containerfile --tag $(application_image_name):latest .

# Runs the project container assuming the image is already built.
container:
	@echo "+$@"
	@echo "############### Removing old container ################"
	@docker rm -f $(application_container_name)

	@echo "################ Running new container ################"
	@docker run --name $(application_container_name) --detach --net host --restart unless-stopped \
		--volume configs/configs.yaml:/etc/$(application_name)/configs.yaml \
		$(application_image_name):latest

# Generates code using the found protocol buffer files.
proto:
	@echo "+$@"
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(proto_path)
