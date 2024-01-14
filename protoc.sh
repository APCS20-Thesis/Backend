protoc -I ./third_party/googleapis \
--proto_path=./api --go_out=api --go_opt=paths=source_relative \
--go-grpc_out=api --go-grpc_opt=paths=source_relative \
--grpc-gateway_out=api --grpc-gateway_opt=paths=source_relative \
api/*.proto