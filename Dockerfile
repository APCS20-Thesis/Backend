# syntax=docker/dockerfile:1

FROM golang:1.20

# Set destination for COPY
WORKDIR /app

RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN protoc --version

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/gogo/protobuf/protoc-gen-gofast@v1.3.1
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.7

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . ./

# Build
# RUN export GOOS=linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o serverrun cmd/*.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 11080

# Run
# CMD ls .
CMD ["./serverrun", "server"]
