package db


import (
	"gopkg.in/mgo.v2/bson"
)


type Community struct{
    Id bson.ObjectId `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Picture string `json:"picture"` // Community picture (header picture?)
    Mods []bson.ObjectId `json:"mods"` // Community Moderators
}


func AddCommunity(u Community) (Community, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("community").Insert(u)

    return u, err
}

func UpdateCommunity(id, newName, newPicture string) (error){
    db := GetDB()
    
    err := db.C("community").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{ "$set": bson.M{"picture": newPicture, "name": newName} })

    return err
}


func DeleteCommunity(id string) (error){
    db := GetDB()

    _, err := db.C("community").RemoveAll(bson.M{"id":id})

    return err
}


func AddModsToCommunity(id, uid string) error{
    db := GetDB()

    err := db.C("community").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$push": bson.M{"mods": bson.ObjectIdHex(uid)}})

    return err
}

func DeleteModsToCommunity(id, uid string) error{
    db := GetDB()

    err := db.C("community").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$pull": bson.M{"mods": bson.ObjectIdHex(uid)}})

    return err
}


