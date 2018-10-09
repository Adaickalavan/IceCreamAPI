package main

import (
	"database"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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

var dictionary = database.Dictionary{}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func main() {

	//Connect to database
	dictionary.Server = os.Getenv("MONGO_PORT")
	dictionary.DatabaseName = os.Getenv("DATABASE_NAME")
	dictionary.CollectionName = os.Getenv("COLLECTION_NAME")
	dictionary.Session = dictionary.Connect()
	defer dictionary.Session.Close()
	//Ensure database index is unique
	dictionary.EnsureIndex([]string{"value"})

	if err := run(); err != nil {
		log.Fatal(err.Error())
	}

}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("LISTENINGADDR")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
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
