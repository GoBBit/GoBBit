package db


import (
	"gopkg.in/mgo.v2/bson"
)

type Poll struct{
    Id bson.ObjectId `json:"id"`
	Title string `json:"title"`
	Options []PollOption `json:"options"` // params for the notification, example: in a mention [topicid, user who mention, mentioned user]
    Tid bson.ObjectId `json:"tid"` // topic ID of the Poll
    Creation_Date int64 `json:"creation_date"`
    Allowed_Options int64 // number of allowed options to select for each user
}

type PollOption struct{
    Title string
    Votes []bson.ObjectId // user ids which voted for this option
}


func AddPoll(u Poll) (Poll, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("poll").Insert(u)

    return u, err
}

func GetPollById(id string) (Poll, error){
    db := GetDB()
    
    u := Poll{}
    err := db.C("poll").Find(bson.M{"id":bson.ObjectIdHex(id)}).One(&u)

    return u, err
}

func GetPollByTopic(id string) (Poll, error){
    db := GetDB()
    
    u := Poll{}
    err := db.C("poll").Find(bson.M{"tid":bson.ObjectIdHex(id)}).One(&u)

    return u, err
}


func DeletePoll(id string) (error){
    db := GetDB()

    _, err := db.C("poll").RemoveAll(bson.M{"id":bson.ObjectIdHex(id)})

    return err
}

func AddVoteToPoll(id, optiontitle, uid string) error{
    db := GetDB()

    err := db.C("poll").Update(bson.M{"id": bson.ObjectIdHex(id), "options.title": optiontitle}, bson.M{"$push": bson.M{"options.$.votes": bson.ObjectIdHex(uid)}})

    return err
}

func DeleteVoteToPoll(id, optiontitle, uid string) error{
    db := GetDB()

    err := db.C("poll").Update(bson.M{"id": bson.ObjectIdHex(id), "options.title": optiontitle}, bson.M{"$pull": bson.M{"options.$.votes": bson.ObjectIdHex(uid)}})

    return err
}

