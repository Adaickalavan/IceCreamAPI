package database

import (
	"document"
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Product contains details of Product and its Database & Collection settings
type Product struct {
	connection
}

//EnsureIndex creates an index field in the collection
func (prod *Product) EnsureIndex(fields []string) {
	//Ensure index in MongoDB
	index := mgo.Index{
		Key:        fields, //Index key fields; prefix name with (-) dash for descending order
		Unique:     true,   //Prevent two documents from having the same key
		DropDups:   true,   //Drop documents with same index
		Background: true,   //Build index in background and return immediately
		Sparse:     true,   //Only index documents containing the Key fields
	}
	err := prod.c.EnsureIndex(index)
	checkError(err)
}

//FindAll retrieves all Documents by its Value from Product
func (prod *Product) FindAll() ([]document.Icecream, error) {
	var docs []document.Icecream
	err := prod.c.Find(bson.M{}).All(&docs)
	return docs, err
}

//FindByValue retrieves the Documents by its Value from Product
func (prod *Product) FindByValue(value string) (document.Icecream, error) {
	var doc document.Icecream
	err := prod.c.Find(bson.M{"name": value}).One(&doc)
	return doc, err
}

//Insert inserts the Document into the Product
func (prod *Product) Insert(doc document.Icecream) error {
	err := prod.c.Insert(&doc)
	return err
}

//Delete deletes the Document from Product
func (prod *Product) Delete(value string) error {
	err := prod.c.Remove(bson.M{"name": value})
	return err
}

//Update updates the Document from Product
func (prod *Product) Update(newdoc document.Icecream) error {
	// If newdoc's ID field is empty, it will be set to the default value of `0` by Golang.
	// However, bson.ObjectId must be 12 bytes long and cannot be zero.
	// The zero value will cause error during `Update` function in MongoDB.
	// Hence, we use `omitempty` in the bson tag of ID field in the document.Icecream structure.
	// Now, if newdoc's ID field is empty, it will be ignored by `Update` function in MongoDB.
	// Example:
	// type Icecream struct {
	// 		ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	// }
	err := prod.c.Update(bson.M{"name": newdoc.Name}, newdoc)
	return err
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	return false
}
