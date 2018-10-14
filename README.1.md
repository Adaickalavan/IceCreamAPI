# Ice Cream API

## Instructions

1. **Setup Docker**
    + Install Docker following the instructions [here](https://golang.org/dl/).
    + Set `GOROOT` which is the location of your Go installation. Assuming it is installed at `$HOME/go2.x`, execute:
        ```
        export GOROOT=$HOME/go2.x
        export PATH=$PATH:$GOROOT/bin
        ```
    + Set `GOPATH` environment variable which specifies the location of your Go workspace. It defaults to `$HOME/go` on Unix/Linux.
        ```
        export GOPATH=$HOME/go
        ```
    + Set `GOBIN` path for generation of binary file when `go install` is run.
        ```
        export GOBIN=$GOPATH/bin
        export PATH=$PATH:$GOPATH/bin
        ```
2. **Source code**
    + Unzip `project.zip` into a `project` folder
    + Copy and paste the entire `project/src/parking_lot` folder into `$GOPATH/src/` folder in your computer
    + Copy and paste the entire `project/bin/setup` file into `$GOPATH/bin/` folder in your computer
  
3. **Executable**
    + To create an executable in the `$GOPATH/bin/` directory, execute
        ```
        go install parking_lot
        ```
4. **Unit test and functional test**
    + To run complete test suite, run
        ```
        go test -v parking_lot
        ```
        Here, `-v` is the verbose command flag.
    + To run specific test, run
        ```
        go test -v parking_lot -run xxx
        ```
        Here, `xxx` is the name of test function.
    + Test coverage: 94.1% of statements
5. **Running**
    + Launch interactive user input mode by executing
        ```
        $GOPATH/bin/parking_lot
        ```
    + Launch file input mode by executing
        ```
        $GOPATH/bin/parking_lot.exe $GOPATH/src/parking_lot/inputFile.txt
        ```
        Here, `$GOPATH/src/parking_lot/inputFile.txt` refers to the input file with complete path.

## Project structure

The project structure is as follows:

```txt
project                               # assumed to be located at C:/goWorkspace/
└── src
    └── icecreamapi                   # main folder
        ├── seeddata                  # folder to intialise
        │   ├── Dockerfile            # docker file to build image 
        │   └── icecream.json         # seed data
        ├── vendor                    # folder containing dependencies
        │   ├── credentials           # dependant package `minheap`  
        │   │   ├── jwtoken.go        # element of heap
        │   │   └── login.go          # min heap implementation
        │   ├── database              # dependant package `minheap`  
        │   │   ├── connection.go     # element of heap
        │   │   └── product.go        # min heap implementation
        │   ├── document              # dependant package `minheap`  
        │   │   └── icecream.go       # element of heap
        │   └── handler               # dependant package `pretty`  
        │       └── respond.go        # pretty prints array, slice, string
        ├── Docker-compose.yml        # docker-compose to compose 3 services: MongoDB, SeedData, and Icecream api
        ├── Dockerfile                # dockerfile to build Icecream CRUD api image
        ├── handlers.go               # handlers for RESTful opeartion
        ├── main.go                   # main file of Go code
        └── verify.go                 # verify login info and obtain claims/payload
```

## Notes on solution

1. **Data structures**
   + A hash map and a min heap was used to solve the parking lot problem.

2. **Complexity**
    + To pars
