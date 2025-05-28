DOCKER_USERNAME=niknike
IMAGE_NAME=bazaar
TAG=latest
APP_NAME=app

run:
	go run ./cmd/app

test:
	mkdir -p coverage
	- go test -v $$(go list ./... | grep -Ev '/(mocks|docs|cmd|db|config|internal/app|generated|internal/transport/dto)') -coverprofile=coverage/cover.out

coverage: test
	go tool cover -html=coverage/cover.out -o coverage/cover.html
	go tool cover -func=coverage/cover.out | grep total:

clean-coverage:
	rm -rf coverage/

mocks:
	go generate ./...

build:
	mkdir -p bin
	go build -o bin/app ./cmd/app

run-build:
	./bin/app

clean:
	rm -rf bin/


.PHONY: docker-build docker-push

# Сборка Docker-образа
docker-build:
	docker build -t $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG) .

# Пуш Docker-образа на Docker Hub
docker-push: docker-build
	docker login
	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)

migrations:
	go run ./cmd/migrations/main.go

.PHONY: swag

swag:
	swag fmt
	swag init -g ./cmd/${APP_NAME}/main.go

auth_proto:
	@mkdir -p internal/transport/generated/auth && \
	protoc --proto_path=proto \
		--go_out=internal/transport/generated/auth \
		--go-grpc_out=internal/transport/generated/auth \
		--go-grpc_opt=paths=source_relative \
		--go_opt=paths=source_relative \
		proto/auth.proto

user_proto:
	@mkdir -p internal/transport/generated/user && \
	protoc --proto_path=proto \
		--go_out=internal/transport/generated/user \
		--go-grpc_out=internal/transport/generated/user \
		--go-grpc_opt=paths=source_relative \
		--go_opt=paths=source_relative \
		proto/user.proto

csat_proto:
	@mkdir -p internal/transport/generated/csat && \
	protoc --proto_path=proto \
		--go_out=internal/transport/generated/csat \
		--go-grpc_out=internal/transport/generated/csat \
		--go-grpc_opt=paths=source_relative \
		--go_opt=paths=source_relative \
		proto/csat.proto

review_proto:
	@mkdir -p internal/transport/generated/review && \
	protoc --proto_path=proto \
		--go_out=internal/transport/generated/review \
		--go-grpc_out=internal/transport/generated/review \
		--go-grpc_opt=paths=source_relative \
		--go_opt=paths=source_relative \
		proto/review.proto

gen-easyjson:
	easyjson -all internal/transport/dto/address.go
	easyjson -all internal/transport/dto/admin.go
	easyjson -all internal/transport/dto/auth.go
	easyjson -all internal/transport/dto/basket.go
	easyjson -all internal/transport/dto/category.go
	easyjson -all internal/transport/dto/csat.go
	easyjson -all internal/transport/dto/minio.go
	easyjson -all internal/transport/dto/notification.go
	easyjson -all internal/transport/dto/order.go
	easyjson -all internal/transport/dto/product.go
	easyjson -all internal/transport/dto/promo.go
	easyjson -all internal/transport/dto/review.go
	easyjson -all internal/transport/dto/search.go
	easyjson -all internal/transport/dto/suggestion.go
	easyjson -all internal/transport/dto/user.go
