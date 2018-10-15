# IceCream API

## Project Summary

A JSON Web Token (JWT) authenticated REST API, to perform CRUD operations on a pre-seeded MongoDB of various icecream products, is developed. The code is completely containerized using Docker images.

## Instructions

1. **Setup Docker**
    + Install Docker following the instructions [here](https://docs.docker.com/docker-for-windows/) for Windows or [here](https://docs.docker.com/docker-for-mac/) for Mac.

2. **Source code**
    + Unzip `icecreamapi.zip` source code into a folder. For example, assume it is unzipped into `C:/goWorkspace/src/icecreamapi` folder.

3. **Build 2 Docker Images** *(requires internet connectivity)*
    + In a bash terminal, navigate to the project folder, i.e., `C:/goWorkspace/src/icecreamapi`.
    + Build a docker image of `icecream` from its Dockerfile by excuting:
        ```bash
        docker build -t "icecream" .
        ```
    + Navigate to the `seeddata` subdirectory in the project folder, i.e. `C:/goWorkspace/src/icecreamapi/seeddata`, by executing:
        ```bash
        cd seeddata/
        ```
    + Build a docker image of `seeddata` from its Dockerfile by excuting:
        ```bash
        docker build -t "seeddata" .
        ```

4. **Run Docker-Compose to Start IceCream API** *(requires internet connectivity)*
    + In a bash terminal, navigate to the project folder, i.e. `C:/goWorkspace/src/icecreamapi`. Start the API by running docker-compose
        ```bash
        docker-compose up
        ```

5. **Functional API Testing**
    + Use an API development environment tool such as [Postman](https://www.getpostman.com) to test the RESTful API endpoints.
    + LOGIN into the `icecream` application by posting `name` and `password` in Postman. Currently, the API only has one authorized user with the below authentication details.
        ```json
        POST http://localhost:8080/login
        BODY
        {
            "name"  : "user1",
            "password" : "1234"
        }
        ```
        A JWT with a validity for 5 minutes, will be returned in the http response, such as:
        ```json
        {
            "tokenString" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzkzMzYwNDEsImlzcyI6IkhvbWVCYXNlIn0.IRKNoD-oEM6j_oDt9DmURkvxYBXmVv76Nlmz_lnoMkU"
        }
        ```
        Please do not use the above token string in your testing as it has exceeded its 5 minutes validity period.
    + Copy and paste the JWT string into the `Token` field under `Auth`->`Bearer Token` in Postman, to be used for subsequent operations.
    + POST a new product into the `icecream` database. New products are posted by `name` and thus the field `name` is mandatory.
        ```json
        POST http://localhost:8080/product
        BODY
        {
            "name"                   : "ChocolateLava",
            "image_closed"           : "",
            "image_open"             : "",
            "description"            : "Molten chocolate",
            "story"                  : "Erupting volcano",
            "sourcing_values"        : ["Non-GMO"],
            "ingredients"            : ["koko","sugar"],
            "allergy_info"           : "none",
            "dietary_certifications" : "Halal",
            "productID"              : "9870"
        }
        ```
        The database enforces unique product `name` and `productID` indexes. Hence, products with same `name` and/or `productID` cannot be inserted into the database.
    + UPDATE a product in the database. Products are updated by `name`, thus the field `name` is mandatory and must match an existing product in the database.
        ```json
        PUT http://localhost:8080/product
        BODY
        {
            "name"                   : "ChocolateLava",
            "image_closed"           : "",
            "image_open"             : "",
            "description"            : "Molten chocolate",
            "story"                  : "Erupting volcano",
            "sourcing_values"        : ["Non-GMO"],
            "ingredients"            : ["koko","sugar, milk, eggs"],
            "allergy_info"           : "contains milk, eggs",
            "dietary_certifications" : "Halal",
            "productID"              : "9870"
        }
        ```
        Upon successfully updating the database, a http response is received as:
        ```json
        {
            "Result": "Successfully updated"
        }
        ```
    + GET a single product from the database. Products are retrieved by `name`, thus the field `name` is mandatory and must match an existing product in the database.
        ```json
        GET http://localhost:8080/product/?doc=ChocolateLava
        ```
        The retrieved product is returned in the http response as:
        ```json
        {
            "name"                   : "ChocolateLava",
            "image_closed"           : "",
            "image_open"             : "",
            "description"            : "Molten chocolate",
            "story"                  : "Erupting volcano",
            "sourcing_values"        : ["Non-GMO"],
            "ingredients"            : ["koko","sugar, milk, eggs"],
            "allergy_info"           : "contains milk, eggs",
            "dietary_certifications" : "Halal",
            "productID"              : "9870"
        }
        ```
    + DELETE a single product from the database. Products are deleted by `name`, thus the field `name` is mandatory and must match an existing product in the database.
        ```json
        DELETE http://localhost:8080/product/?doc=ChocolateLava
        ```
        Upon successfully deleting the product, a http response is received as:
        ```json
        {
            "Result": "Successfully deleted"
        }
        ```
    + GET all products from the database.
        ```json
        GET http://localhost:8080/product/
        ```
        All products from the databse is returned in the http response in JSON format.

## Project Structure

The project structure is as follows:

```txt
project                               # assumed to be located at C:/goWorkspace/
└── src                               #
    └── icecreamapi                   # main folder
        ├── seeddata                  # special folder to import initial data into MongoDB
        │   ├── Dockerfile            # dockerfile for building image to seed data into MongoDB
        │   └── icecream.json         # seed data
        ├── vendor                    # folder containing dependencies
        │   ├── credentials           # dependant package `credentials`  
        │   │   ├── jwtoken.go        # create and authenticate JWT
        │   │   └── login.go          # hash and compare login passwords
        │   ├── database              # dependant package `database`  
        │   │   ├── connection.go     # generic, reusable, database connection function
        │   │   └── product.go        # database CRUD operations
        │   ├── document              # dependant package `document`  
        │   │   └── icecream.go       # define the `icecream` document to be stored in the MongoDB
        │   └── handler               # dependant package `handler`  
        │       └── respond.go        # generic, reusable, http response functions
        ├── Docker-compose.yml        # to compose 3 services: `mongo`, `seeddata`, and `icecream`
        ├── Dockerfile                # dockerfile to build `icecream` api image
        ├── handlers.go               # handlers for RESTful operation
        ├── main.go                   # main file of Go code
        └── verify.go                 # verify user login info and obtain claims/payload
```

## Notes on Solution

1. **Language and Structure**
   + The code is written in Golang.
   + MongoDB is used as the store database.
   + A server with REST endpoints listens at port `localhost:8080`.
   + The code is completely containerized for easy deployment, with 3 docker images. Namely,
        + `mongo` - to create Mongo database. This image will be directly pulled via internet from the Docker Hub,
        + `seeddata` - to initialize the database with seed data, and
        + `icecream` - for CRUD operations.

2. **Authentication**
    + Only authorized users are able to access the REST endpoints to perform CRUD operations on the database.
    + Users are authenticated using JWT mechanism, as it is advantageous for scalability.
    + Currently, only one authenticated user is present: `{"name":"user1","password":"1234"}`. Passwords stored in the source code are hashed using bcrypt algorithm.
    + Additionally, MongoDB access is protected with an username and password (`--username admin1 --password abcd --authenticationDatabase admin`).

3. **Docker Notes**
   + When docker-compose is run with Docker-Toolbox, go to `192.168.99.100:8080/` to interact with the application. `192.168.99.100` is the IP address of your docker-machine. Execute `$ docker-machine ip` to get IP address of your docker-machine.
   + To tear down all containers and stored volume: `docker-compose down -v`
   + To prune all dangling containers, networks, and build caches: `docker system prune`