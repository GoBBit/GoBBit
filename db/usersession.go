package db


import (
    "gopkg.in/mgo.v2/bson"

)


type UserSession struct{
    Id string `json:"id",omitempty`
    Uid bson.ObjectId `json:"uid",omitempty`
}


func AddUserSession(u UserSession) (UserSession, error){
    db := GetDB()

    err := db.C("session").Insert(u)

    return u, err
}

func DeleteUserSession(id string) (error){
    db := GetDB()

    _, err := db.C("session").RemoveAll(bson.M{"id":id})

    return err
}

func GetUserSession(id string) (UserSession, error){
    db := GetDB()
    
    u := UserSession{}
    err := db.C("session").Find(bson.M{"id":id}).One(&u)

    return u, err
}