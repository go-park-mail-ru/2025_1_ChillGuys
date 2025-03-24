DOCKER_USERNAME=niknike
IMAGE_NAME=bazaar
TAG=latest
APP_NAME=app

run:
	go run ./cmd/app

test:
	mkdir -p coverage
	go test -v $$(go list ./... | grep -v '/mocks') -coverprofile=coverage/cover.out

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
