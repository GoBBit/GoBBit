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
	Community string `json:"community"` // community slug
	Creation_Date int64 `json:"creation_date"`
    Editation_Date int64 `json:"editation_date"`
    Last_Update int64 `json:"last_update"`
    Posts_number int64 `json:"posts_number"` // number of posts(messages) on the topic
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

func UpdateTopicLastUpdate(id string, t int64) error{
    db := GetDB()

    err := db.C("topic").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"last_update": t}})

    return err
}

func GetTopicsByCommunity(cslugs []string, limit, start int) ([]Topic, error){
    db := GetDB()
    
    u := []Topic{}
    err := db.C("topic").Find(bson.M{"community": bson.M{"$in": cslugs }}).Skip(start).Limit(limit).Sort("-last_update").All(&u)

    return u, err
}

func GetTopicsByCommunityWithoutIgnoredUsers(cslugs []string, limit, start int, ignored []bson.ObjectId) ([]Topic, error){
    db := GetDB()
    
    u := []Topic{}
    q := bson.M{
        "community": bson.M{"$in": cslugs }, 
        "uid": bson.M{"$nin": ignored },
    }
    err := db.C("topic").Find(q).Skip(start).Limit(limit).Sort("-last_update").All(&u)

    return u, err
}

func GetTopicsListWithoutIgnoredUsers(limit, start int, ignored []bson.ObjectId) ([]Topic, error){
    db := GetDB()
    
    u := []Topic{}
    q := bson.M{
        "uid": bson.M{"$nin": ignored },
    }
    err := db.C("topic").Find(q).Skip(start).Limit(limit).Sort("-last_update").All(&u)

    return u, err
}

func GetTopicsByUser(uid string, limit, start int) ([]Topic, error){
    db := GetDB()
    
    u := []Topic{}
    err := db.C("topic").Find(bson.M{"uid": bson.ObjectIdHex(uid)}).Skip(start).Limit(limit).Sort("-last_update").All(&u)

    return u, err
}

func (t *Topic) IncrementPostsNumber(n int) error{
    db := GetDB()

    err := db.C("topic").Update(bson.M{"id": t.Id}, bson.M{"$inc": bson.M{"posts_number": n}})

    return err
}

func (t *Topic) GenerateSlug() (string){
    t.Slug = slug.Slug(t.Title)
    return t.Slug
}

func (t *Topic) IsOwner(u User) (bool){
    return (t.Uid == u.Id)
}

