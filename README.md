# DictionaryAPI

This project builds a REST web API to POST and GET word defintions to/from a dictionary stored in a MongoDB database. Dockerfile and Docker-compose files are provided to containerize the deployment of Go code and MongoDB database.

Docker commands:

1. Build image of Go code:
   + `docker build -t "dictionaryapi" .`
2. Create and run all containers:
   + `docker-compose up`
3. Tear down all containers and stored volume:
   + `docker-compose down -v`

Further extension of functionality and more description of the code will be proivded later.

See [website](https://adaickalavan.github.io/portfolio/docker_golang_rest_kafka_mongodb/) for information.


Need unrestricted connectivity to 
1. Github
2. Docker Hub

remember that product id must be unique or it will not be inserted

database is unique by name and productid

delete is done through name
update is done trhough name
no name will result in error while upadting

We will be using the bcrypt algorithm to hash and salt our passwords.

for automated docker build
