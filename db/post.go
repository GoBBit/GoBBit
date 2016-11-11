package db


import (
	"gopkg.in/mgo.v2/bson"

)


type Post struct{
    Id bson.ObjectId `json:"id"`
	Uid bson.ObjectId `json:"uid"`
	Tid bson.ObjectId `json:"tid"` // topic ID
	Content string `json:"content"`
    Creation_Date int64 `json:"creation_date"`
    Editation_Date int64 `json:"editation_date"`
}

func AddPost(u Post) (Post, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("post").Insert(u)

    return u, err
}

func UpdatePost(id, newContent string) (error){
    db := GetDB()
    
    err := db.C("post").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{ "$set": bson.M{"content": newContent} })

    return err
}


func DeletePost(id string) (error){
    db := GetDB()

    _, err := db.C("post").RemoveAll(bson.M{"id":id})

    return err
}


