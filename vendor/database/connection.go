package database

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

//Connection contains server connection settings
type connection struct {
	Server         string          //Connection to server 'IP:Port'
	DatabaseName   string          //Name of desired database
	CollectionName string          //Name of desired collection
	UserName       string          //Username to login into database
	Password       string          //Password to login into database
	Session        *mgo.Session    //Session
	c              *mgo.Collection //Pointer to desired collection
}

//Connect connects to the database
func (conn *connection) Connect() *mgo.Session {
	info := &mgo.DialInfo{
		Addrs:    []string{conn.Server},
		Timeout:  60 * time.Second,
		Username: conn.UserName,
		Password: conn.Password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal("Database dial error:", err)
	}
	conn.c = session.DB(conn.DatabaseName).C(conn.CollectionName)
	return session
}
