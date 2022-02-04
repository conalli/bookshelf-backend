#! /bin/sh

if [ $1 = "dev" ]
then
  docker-compose -f ./build/docker-compose.dev.yml -p bookshelf_dev build
  exit 0
elif [ $1 = "test" ]
then
  docker-compose -f ./build/docker-compose.test.yml -p bookshelf_test build
  exit 0
elif [ $1 = "local" ]
then
  docker-compose -f ./build/docker-compose.local.yml -p bookshelf_local build
  exit 0
else
  exit 1
fi
