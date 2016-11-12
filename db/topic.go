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
	Community string `json:"cid"` // community slug
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

func UpdateTopic(u Topic) error{
    db := GetDB()

    err := db.C("topic").Update(bson.M{"id": u.Id}, u)

    return err
}

func GetTopicById(id string) (Topic, error){
    db := GetDB()
    
    u := Topic{}
    err := db.C("topic").Find(bson.M{"id":bson.ObjectIdHex(id)}).One(&u)

    return u, err
}


func DeleteTopic(id string) (error){
    db := GetDB()

    _, err := db.C("topic").RemoveAll(bson.M{"id":bson.ObjectIdHex(id)})

    return err
}


func (t *Topic) GenerateSlug() (string){
    t.Slug = slug.Slug(t.Title)
    return t.Slug
}

