#! /bin/sh

if [ $1 = "dev" ]
then
  docker-compose -p bookshelf_dev down -v
  exit 0
elif [ $1 = "test" ]
then
  docker-compose -p bookshelf_test down -v
  exit 0
elif [ $1 = "local" ]
then
  docker-compose -p bookshelf_local down -v
  exit 0
else
  exit 1
fi
