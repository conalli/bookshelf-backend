#! /bin/sh

if [ $1 = "dev" ]
then
  docker-compose -f ./build/docker-compose.dev.yml -p bookshelf_dev up -d
  exit 0
elif [ $1 = "test" ]
then
  docker-compose -f ./build/docker-compose.test.yml -p bookshelf_test up -d
  exit 0
else
  exit 1
fi
