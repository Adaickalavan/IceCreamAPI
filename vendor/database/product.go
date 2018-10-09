package database

import (
	"document"
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Connection contains server and database strings
type Connection struct {
	Server         string          //Connection to server 'IP:Port'
	DatabaseName   string          //Name of desired database
	CollectionName string          //Name of desired collection
	Session        *mgo.Session    //Session
	c              *mgo.Collection //Pointer to desired database
	// db             *mgo.Database   //Pointer to desired database
}

//Connect connects to the database
func (conn Connection) Connect() *mgo.Session {
	info := &mgo.DialInfo{
		Addrs:    []string{conn.Server},
		Timeout:  60 * time.Second,
		Database: conn.DatabaseName,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal("Database dial error:", err)
	}
	conn.c = session.DB(conn.DatabaseName).C(conn.CollectionName)
	return session
}

//EnsureIndex creates an index field in the collection
func (conn *Connection) EnsureIndex(fields []string) {
	//Ensure index in MongoDB
	index := mgo.Index{
		Key:        fields, //Index key fields; prefix name with (-) dash for descending order
		Unique:     true,   //Prevent two documents from having the same key
		DropDups:   true,   //Drop documents with same index
		Background: true,   //Build index in background and return immediately
		Sparse:     true,   //Only index documents containing the Key fields
	}
	err := conn.c.EnsureIndex(index)
	checkError(err)
}

//FindAll retrieves all Documents by its Value from Connection
func (conn *Connection) FindAll() ([]document.Icecream, error) {
	var docs []document.Icecream
	err := conn.c.Find(bson.M{}).All(&docs)
	return docs, err
}

//FindByValue retrieves the Documents by its Value from Connection
func (conn *Connection) FindByValue(value string) (document.Icecream, error) {
	var doc document.Icecream
	err := conn.c.Find(bson.M{"value": value}).One(&doc)
	return doc, err
}

//Insert inserts the Document into the Connection
func (conn *Connection) Insert(doc document.Icecream) error {
	err := conn.c.Insert(&doc)
	return err
}

//Delete deletes the doc from Connection
func (conn *Connection) Delete(doc document.Icecream) error {
	err := conn.c.Remove(&doc)
	return err
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	return false
}
