LOCAL_BIN:=$(CURDIR)\bin

install-deps:
	set GOBIN=$(LOCAL_BIN) && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	set GOBIN=$(LOCAL_BIN) && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-auth-api

generate-auth-api:
	mkdir pkg\auth_v1 || exit 0
	protoc --proto_path=grpc\api\auth_v1 --go_out=pkg\auth_v1 --go_opt=paths=source_relative --plugin=bin\protoc-gen-go --go-grpc_out=pkg\auth_v1 --go-grpc_opt=paths=source_relative --plugin=bin\protoc-gen-go-grpc grpc\api\auth_v1\auth.proto

	--go_out=pkg\auth_api_v1 --go_opt=paths=source_relative ^
	--plugin=protoc-gen-go=bin\protoc-gen-go ^
	--go-grpc_out=pkg\auth_api_v1 --go-grpc_opt=paths=source_relative ^
	--plugin=protoc-gen-go-grpc=bin\protoc-gen-go-grpc ^
	grpc\api\auth_v1\auth.proto