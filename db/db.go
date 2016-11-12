package db


import (
    //"log"
    //"time"
    "os"
    //"io"
    //"strings"
    //"fmt"
    //"net/http"
    //"html"
    //"encoding/json"

    "gopkg.in/mgo.v2"

)


type DB struct{
	MongoSession *mgo.Session
	Name string
	Host string
}

var DB_HOST = map[bool]string{true: os.Getenv("DB_HOST"), false: "localhost"} [os.Getenv("DB_HOST") != ""]
var DB_USER = map[bool]string{true: os.Getenv("DB_USER"), false: ""} [os.Getenv("DB_USER") != ""]
var DB_PASS = map[bool]string{true: os.Getenv("DB_PASS"), false: ""} [os.Getenv("DB_PASS") != ""]

var DB_NAME = map[bool]string{true: os.Getenv("DB_NAME"), false: "GoBBit"} [os.Getenv("DB_NAME") != ""]

var instance *DB = nil

func GetInstance() *DB {
    if instance == nil {
        instance = &DB{Name:DB_NAME, Host:DB_HOST}

        sess, err := mgo.Dial(instance.Host)
        if err != nil{
        	panic(err)
        }

        instance.MongoSession = sess

        if DB_USER != "" && DB_PASS != ""{
        	instance.MongoSession.DB(instance.Name).Login(DB_USER, DB_PASS)
        	
        	if err != nil{
        		panic(err)
        	}
        }

    }
    return instance
}

func GetDB() (*mgo.Database){
	db := GetInstance()
	return db.MongoSession.DB(db.Name)
}

func EnsureIndex(){
    db := GetDB()
    
    db.C("user").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("user").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true, DropDups: true})
    db.C("user").EnsureIndex(mgo.Index{Key: []string{"slug"}, Unique: true, DropDups: true})

    db.C("session").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("session").EnsureIndex(mgo.Index{Key: []string{"uid"}})

    db.C("community").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("community").EnsureIndex(mgo.Index{Key: []string{"slug"}, Unique: true, DropDups: true})

    db.C("topic").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("topic").EnsureIndex(mgo.Index{Key: []string{"community"}, Unique: false})
    db.C("topic").EnsureIndex(mgo.Index{Key: []string{"last_update"}, Unique: false})

    db.C("post").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("post").EnsureIndex(mgo.Index{Key: []string{"creation_date"}, Unique: true, DropDups: true})

}



