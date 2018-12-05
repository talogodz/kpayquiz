package database

import (
	"log"

	mgo "github.com/globalsign/mgo"
)

type DAO struct {
	Server   string
	Database string
}

var (
	DB *mgo.Database
)

func (db *DAO) Connect() {
	session, err := mgo.Dial(db.Server)
	if err != nil {
		log.Fatal(err)
	}
	DB = session.DB(db.Database)
}
