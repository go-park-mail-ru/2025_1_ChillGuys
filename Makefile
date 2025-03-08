test:
	mkdir -p coverage
	go test -v $$(go list ./... | grep -v '/mocks') -coverprofile=coverage/cover.out

coverage: test
	go tool cover -html=coverage/cover.out -o coverage/cover.html
	go tool cover -func=coverage/cover.out | grep total:

mocks:
	cd ./internal/transport && go generate

build:
	mkdir -p bin
	go build -o bin/app ./cmd/app

run: build
	./bin/app

#.PHONY: docker-build docker-push
#
#DOCKER_USERNAME=
#IMAGE_NAME=
#TAG=latest=
#
## Сборка Docker-образа
#docker-build:
#	docker build -t $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG) ./.
#
## Пуш Docker-образа на Docker Hub
#docker-push: docker-build
#	docker login
#	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)
#
