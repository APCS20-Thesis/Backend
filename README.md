# CDP Backend

### Set up (Mac)
##### Golang
Install Golang, choose compatible package https://go.dev/dl/.

Check if _the package installs the Go distribution to `/usr/local/go`_

```shell
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
##### Protobuf
```shell
brew install protobuf
```

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 
go install github.com/gogo/protobuf/protoc-gen-gofast@v1.3.1
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.7
```

##### MigrateCLI
```shell
brew install golang-migrate
```

### Instructions

##### Run server
Postgres DB is included in `docker-compose.yml`. Run DB:
```shell
docker compose up
```
Run server
```shell
go run cmd/main.go server
```

##### Implement new API
1. Define API schema in `api/api.proto` and `api/data.proto`
2. Run `protoc.sh` to generate API 
3. Implement API in `internal/service`

##### Data Migrations
1. Create new migration file
```shell 
migrate create -dir ./sql/migrations -ext sql "write_description_for_migration" 
```
2. Write migrations in SQL in that two new files (.up.sql and .down.sql)
3. To migrate to the latest migration, run 
```shell
go run cmd/main.go migrate up
```
4. To revert the migrations, run 
```shell
go run cmd/main.go migrate down NUMBER_OF_VERSIONS
```

