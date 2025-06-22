# --- ПЕРЕМЕННЫЕ ---
include .env
LOCAL_BIN=$(CURDIR)/bin

# --- УСТАНОВКА ---
install-deps:
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2

# --- ПРОТО ---
vendor-proto:
	@rm -rf vendor.protogen/googleapis
	@git clone --depth=1 https://github.com/googleapis/googleapis vendor.protogen/googleapis

generate-auth-api:
	@mkdir -p pkg/auth_v1
	@PATH=$(LOCAL_BIN):$$PATH \
	protoc \
		--proto_path=api/auth_v1 \
		--proto_path=vendor.protogen/googleapis \
		--proto_path=. \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pkg/auth_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		api/auth_v1/auth.proto

re-generate-auth-api:
	@rm -f pkg/auth_v1/auth*.pb.go
	@$(MAKE) generate-auth-api

# --- ПРОЧЕЕ ---
get-deps:
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	@go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@109.71.15.36:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/beachrockhotel/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAAkL-JL7tX1UmJizGV3dsIj9cYqY7Y0WEq cr.selcloud.ru/beachrockhotel
	docker push cr.selcloud.ru/beachrockhotel/test-server:v0.0.1

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
