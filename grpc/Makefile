LOCAL_BIN=$(CURDIR)/bin

install-deps:
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	@go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

generate:
	make generate-auth-api

copy-to-server:
	scp service_linux root@109.71.15.36:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/beachrockhotel/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAA4K83Tm9QTg-rGLukLjESun0YL1NQVIsR cr.selcloud.ru/beachrockhotel
	docker push cr.selcloud.ru/beachrockhotel/test-server:v0.0.1

generate-auth-api:
	@mkdir -p pkg/auth_v1
	@protoc \
		--proto_path=api/auth_v1 \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		api/auth_v1/auth.proto