package main

import (
	"database"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

//Hooks that may be overridden for testing
var inputReader io.Reader = os.Stdin
var outputWriter io.Writer = os.Stdout

func init() {
	//Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

}

var product = database.Product{}

func main() {

	//Connect to database
	product.Server = os.Getenv("MONGO_PORT")
	product.DatabaseName = os.Getenv("DATABASE_NAME")
	product.CollectionName = os.Getenv("COLLECTION_NAME")
	product.Session = product.Connect()
	defer product.Session.Close()

	//Ensure database index is unique
	product.EnsureIndex([]string{"productID"})
	product.EnsureIndex([]string{"name"})

	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	httpAddr := os.Getenv("LISTENING_ADDR")
	mux := makeMuxRouter()
	loggedRouter := handlers.LoggingHandler(outputWriter, mux) //Wrap the mux router to log all api requests. Logged requests are written to `outputWriter`
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        loggedRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening on ", httpAddr)
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func checkError(err error) bool {
	if err != nil {
		fmt.Fprintln(outputWriter, err.Error())
		return true
	}
	return false
}
