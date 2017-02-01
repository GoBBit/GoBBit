package db


import (
	"gopkg.in/mgo.v2/bson"
)

// Update it when new notifications are added or when params are modified
// Types and params:
//  - mention
//      [tid, user_slug, mentioned_user_slug]
//  - new_post (for topic creator)
//      [tid, user_slug]


type Notification struct{
    Id bson.ObjectId `json:"id"`
	Type string `json:"type"` // type of notification (mention, new_post...)
	Params []string `json:"params"` // params for the notification, example: in a mention [topicid, user who mention, mentioned user]
    Uid bson.ObjectId `json:"uid"` // sent to User ID
    Read bool `json:"read"`
}


func AddNotification(u Notification) (Notification, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("notification").Insert(u)

    return u, err
}

func GetNotificationsByUser(uid string) (Notification, error){
    db := GetDB()
    
    u := Notification{}
    err := db.C("notification").Find(bson.M{"uid":bson.ObjectIdHex(uid)}).One(&u)

    return u, err
}

func DeleteNotification(id string) (error){
    db := GetDB()

    _, err := db.C("notification").RemoveAll(bson.M{"id":bson.ObjectIdHex(id)})

    return err
}

func DeleteNotificationsByUser(uid string) (error){
    db := GetDB()

    _, err := db.C("notification").RemoveAll(bson.M{"uid":bson.ObjectIdHex(uid)})

    return err
}

func MarkAsReadAllNotificationsByUser(uid string) (error){
    db := GetDB()

    err := db.C("notification").Update(bson.M{"uid": bson.ObjectIdHex(uid)}, bson.M{"$set": bson.M{"read": true}})

    return err
}

func (n *Notification) MarkAsRead() (){
    db := GetDB()
    n.Read = true
    db.C("notification").Update(bson.M{"id": n.Id}, bson.M{"$set": bson.M{"read": true}})
}

