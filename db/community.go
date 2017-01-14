package db


import (
	"gopkg.in/mgo.v2/bson"
    "github.com/tv42/slug"

    "GoBBit/utils"
)


type Community struct{
    Id bson.ObjectId `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
    Description string `json:"description"`
	Picture string `json:"picture"` // Community picture (header picture?)
    Mods []bson.ObjectId `json:"mods"` // Community Moderators
    Banned_Users []bson.ObjectId `json:"banned_users"` // Community Moderators
    Creation_Date int64 `json:"creation_date"`
    Posts_number int64 `json:"posts_number"` // number of total posts(messages) on the community
    Topics_number int64 `json:"topics_number"` // number of total topics on the community
}


func AddCommunity(u Community) (Community, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("community").Insert(u)

    return u, err
}

func UpdateCommunity(u Community) (error){
    db := GetDB()
    
    err := db.C("community").Update(bson.M{"id": u.Id}, u)

    return err
}

func GetCommunityBySlug(slug string) (Community, error){
    db := GetDB()
    
    u := Community{}
    err := db.C("community").Find(bson.M{"slug":slug}).One(&u)

    return u, err
}

func DeleteCommunity(id string) (error){
    db := GetDB()

    _, err := db.C("community").RemoveAll(bson.M{"id":bson.ObjectIdHex(id)})

    return err
}

func DeleteCommunityBySlug(slug string) (error){
    db := GetDB()

    _, err := db.C("community").RemoveAll(bson.M{"slug":slug})

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

func AddBannedUserToCommunity(id, uid string) error{
    db := GetDB()

    err := db.C("community").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$push": bson.M{"banned_users": bson.ObjectIdHex(uid)}})

    return err
}

func DeleteBannedUserToCommunity(id, uid string) error{
    db := GetDB()

    err := db.C("community").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$pull": bson.M{"banned_users": bson.ObjectIdHex(uid)}})

    return err
}

func (c *Community) GenerateSlug() (string){
    c.Slug = slug.Slug(c.Name)
    return c.Slug
}

func (c *Community) IsMod(u User) (bool){
    return (utils.IndexOf(c.Mods, u.Id) >= 0)
}

func (c *Community) IsBanned(u User) (bool){
    return (utils.IndexOf(c.Banned_Users, u.Id) >= 0)
}

func (c *Community) UserCanPost(u User) (bool){
    // if not banned in the community nor the forum
    return !c.IsBanned(u) && !u.IsBanned
}

func (c *Community) IncrementPostsNumber(n int) error{
    db := GetDB()

    err := db.C("community").Update(bson.M{"id": c.Id}, bson.M{"$inc": bson.M{"posts_number": n}})

    return err
}

func (c *Community) IncrementTopicsNumber(n int) error{
    db := GetDB()

    err := db.C("community").Update(bson.M{"id": c.Id}, bson.M{"$inc": bson.M{"topics_number": n}})

    return err
}

