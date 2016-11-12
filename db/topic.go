package db


import (
	"gopkg.in/mgo.v2/bson"
    "github.com/tv42/slug"
)


type Topic struct{
    Id bson.ObjectId `json:"id"`
    Title string `json:"title"`
    Slug string `json:"slug"`
    Content string `json:"content"` // main post content
	Uid bson.ObjectId `json:"uid"`
	Cid bson.ObjectId `json:"cid"` // community ID
	Creation_Date int64 `json:"creation_date"`
    Editation_Date int64 `json:"editation_date"`
    Last_Update int64 `json:"last_update"`
}


func AddTopic(u Topic) (Topic, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("topic").Insert(u)

    return u, err
}

func UpdateTopic(id, newTitle, newContent string) (error){
    db := GetDB()

    err := db.C("topic").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{ "$set": bson.M{"content": newContent, "title": newTitle} })

    return err
}


func DeleteTopic(id string) (error){
    db := GetDB()

    _, err := db.C("topic").RemoveAll(bson.M{"id":id})

    return err
}


func (t *Topic) GenerateSlug() (string){
    t.Slug = slug.Slug(t.Title)
    return t.Slug
}

