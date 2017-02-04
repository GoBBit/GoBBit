package db


import (
    "gopkg.in/mgo.v2"

    "GoBBit/config"
)


type DB struct{
	MongoSession *mgo.Session
	Name string
	Host string
}

var instance *DB = nil

func GetInstance() *DB {
    if instance == nil {
        dbConfig := config.GetInstance().DbConfig
        instance = &DB{Name:dbConfig.Name, Host:dbConfig.Host}

        sess, err := mgo.Dial(instance.Host)
        if err != nil{
        	panic(err)
        }

        instance.MongoSession = sess

        if dbConfig.User != "" && dbConfig.Pass != ""{
        	instance.MongoSession.DB(instance.Name).Login(dbConfig.User, dbConfig.Pass)
        	
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

    db.C("community").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("community").EnsureIndex(mgo.Index{Key: []string{"slug"}, Unique: true, DropDups: true})

    db.C("topic").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("topic").EnsureIndex(mgo.Index{Key: []string{"community"}, Unique: false})
    db.C("topic").EnsureIndex(mgo.Index{Key: []string{"last_update"}, Unique: false})

    db.C("post").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("post").EnsureIndex(mgo.Index{Key: []string{"creation_date"}, Unique: false})

    db.C("notification").EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true})
    db.C("notification").EnsureIndex(mgo.Index{Key: []string{"creation_date"}, Unique: false})
    db.C("notification").EnsureIndex(mgo.Index{Key: []string{"uid"}, Unique: false})

}



