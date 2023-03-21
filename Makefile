.PHONY: all build up down

all: build up

build:
	docker-compose -f ./build/docker-compose.dev.yml -p bookshelf_dev build

up:
	docker-compose -f ./build/docker-compose.dev.yml -p bookshelf_dev up -d

down:
	docker-compose -p bookshelf_dev down -v
	rm -rf ./tmp