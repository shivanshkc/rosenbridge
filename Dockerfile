#------------------------------------------------------------------
FROM golang:1.17-alpine as builder

# Update alpine.
RUN apk update && apk upgrade

# Install alpine dependencies.
RUN apk --no-cache --update add build-base bash

# Create and change to the 'project' directory.
WORKDIR /project

# Install project dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy and test code.
COPY . .
RUN make test

# Build application binary.
RUN make build-alpine

#-------------------------------------------------------------------
FROM alpine:3

# Create and change to the the 'project' directory.
WORKDIR /project

# Copy the files to the production image from the builder stage.
COPY --from=builder /project/bin /project/

# Run the web service on container startup.
CMD ["/project/main"]

#-------------------------------------------------------------------
