include .env

LOCAL_BIN=$(CURDIR)/bin

install-deps:
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	@GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	@GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

get-deps:
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	@go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

generate:
	make generate-auth-api
	make generate-access-api

generate-auth-api:
	@mkdir -p pkg/auth_v1
	protoc --proto_path=api/auth_v1 --proto_path=vendor.protogen \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		--validate_out=lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
		--grpc_gateway_out=pkg/auth_v1 --grpc_gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc_gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
		--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
		api/auth_v1/auth.proto
	@test -f pkg/swagger/api.swagger.json || (echo "Swagger JSON not found!"; exit 1)
	@$(MAKE) generate-statik

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/access_v1/access.proto

gen-cert:
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
    		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext

generate-statik:
	@$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png' -f -p statik

re-generate-auth-api:
	@rm -f pkg/auth_v1/auth*.pb.go
	@$(MAKE) generate-auth-api

copy-to-server:
	scp service_linux root@109.71.15.36:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/beachrockhotel/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAAkL-JL7tX1UmJizGV3dsIj9cYqY7Y0WEq cr.selcloud.ru/beachrockhotel
	docker push cr.selcloud.ru/beachrockhotel/test-server:v0.0.1

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" down -v

vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate && \
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate && \
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate && \
		rm -rf vendor.protogen/protoc-gen-validate ; \
	fi

	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		mkdir -p vendor.protogen/google && \
		mv vendor.protogen/googleapis/google/api vendor.protogen/google && \
		rm -rf vendor.protogen/googleapis ; \
	fi

	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi

grpc-load-test:
	cd ghz_tmp && ../bin/ghz \
		--proto auth_v1/auth.proto \
		--call auth_v1.AuthV1.Get \
		--data '{"id": 1}' \
		--rps 100 \
		--total 3000 \
		--cacert ../ca.cert \
		--cert ../service.pem \
		--key ../service.key \
		localhost:50051

grpc-error-load-test:
	cd ghz_tmp && ../bin/ghz \
		--proto auth_v1/auth.proto \
		--call auth_v1.AuthV1.Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--cacert ../ca.cert \
		--cert ../service.pem \
		--key ../service.key \
		localhost:50051
