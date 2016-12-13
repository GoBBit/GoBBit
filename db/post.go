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

func UpdatePost(u Post) (error){
    db := GetDB()
    
    err := db.C("post").Update(bson.M{"id": u.Id}, u)

    return err
}

func GetPostById(id string) (Post, error){
    db := GetDB()
    
    u := Post{}
    err := db.C("post").Find(bson.M{"id":bson.ObjectIdHex(id)}).One(&u)

    return u, err
}

func GetPostsByTopicId(tid string) ([]Post, error){
    db := GetDB()
    
    u := []Post{}
    err := db.C("post").Find(bson.M{"tid":bson.ObjectIdHex(tid)}).All(&u)

    return u, err
}


func DeletePost(id string) (error){
    db := GetDB()

    _, err := db.C("post").RemoveAll(bson.M{"id":bson.ObjectIdHex(id)})

    return err
}


func (p *Post) IsOwner(u User) (bool){
    return (p.Uid == u.Id)
}

