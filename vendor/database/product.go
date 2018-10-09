package database

import (
	"document"
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//COLLECTION is the name of collection within Dictionary database
const COLLECTION = "icecream"

//Dictionary contains server and database strings
type Dictionary struct {
	Server       string        //Connection to server 'IP:Port'
	DatabaseName string        //Name of desired database
	Session      *mgo.Session  //Session
	db           *mgo.Database //Pointer to desired database
}

//Connect connects to the database
func (dictionary Dictionary) Connect() *mgo.Session {
	info := &mgo.DialInfo{
		Addrs:    []string{dictionary.Server},
		Timeout:  60 * time.Second,
		Database: dictionary.DatabaseName,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal("Database dial error:", err)
	}
	dictionary.db = session.DB(dictionary.DatabaseName)
	return session
}

//EnsureIndex creates an index field in the collection
func (dictionary *Dictionary) EnsureIndex(fields []string) {
	//Ensure index in MongoDB
	index := mgo.Index{
		Key:        fields, //Index key fields; prefix name with (-) dash for descending order
		Unique:     true,   //Prevent two documents from having the same key
		DropDups:   true,   //Drop documents with same index
		Background: true,   //Build index in background and return immediately
		Sparse:     true,   //Only index documents containing the Key fields
	}
	err := dictionary.db.C(COLLECTION).EnsureIndex(index)
	checkError(err)
}

//FindAll retrieves all doc by its Value from dictionary
func (dictionary *Dictionary) FindAll() ([]document.Icecream, error) {
	var docs []document.Icecream
	err := dictionary.db.C(COLLECTION).Find(bson.M{}).All(&docs)
	return docs, err
}

//FindByValue retrieves the doc by its Value from dictionary
func (dictionary *Dictionary) FindByValue(value string) (document.Icecream, error) {
	var doc document.Icecream
	err := dictionary.db.C(COLLECTION).Find(bson.M{"value": value}).One(&doc)
	return doc, err
}

//Insert inserts the doc into the dictionary
func (dictionary *Dictionary) Insert(doc document.Icecream) error {
	err := dictionary.db.C(COLLECTION).Insert(&doc)
	return err
}

//Delete deletes the doc from dictionary
func (dictionary *Dictionary) Delete(doc document.Icecream) error {
	err := dictionary.db.C(COLLECTION).Remove(&doc)
	return err
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	return false
}
