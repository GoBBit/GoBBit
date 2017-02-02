package db


import (
	"gopkg.in/mgo.v2/bson"

    "regexp" // for create notifications
    "time"
    "encoding/json"
)

// Update it when new notifications are added or when params are modified
// Types and params:
//  - mention
//      [tid, user_slug]
//  - new_post (for topic creator)
//      [tid, user_slug]


type Notification struct{
    Id bson.ObjectId `json:"id"`
	Type string `json:"type"` // type of notification (mention, new_post...)
	Params []string `json:"params"` // params for the notification, example: in a mention [topicid, user who mention, mentioned user]
    Uid bson.ObjectId `json:"uid"` // sent to User ID
    Read bool `json:"read"`
    Creation_Date int64 `json:"creation_date"`
}


func AddNotification(u Notification) (Notification, error){
    db := GetDB()
    
    u.Id = bson.NewObjectId()
    err := db.C("notification").Insert(u)

    return u, err
}

func GetNotificationById(id string) (Notification, error){
    db := GetDB()
    
    u := Notification{}
    err := db.C("notification").Find(bson.M{"id":bson.ObjectIdHex(id)}).One(&u)

    return u, err
}

func GetNotificationsByUser(uid string) ([]Notification, error){
    db := GetDB()
    
    u := []Notification{}
    err := db.C("notification").Find(bson.M{"uid":bson.ObjectIdHex(uid)}).Sort("-creation_date").All(&u)

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

    _, err := db.C("notification").UpdateAll(bson.M{"uid": bson.ObjectIdHex(uid)}, bson.M{"$set": bson.M{"read": true}})

    return err
}

func (n *Notification) MarkAsRead() (){
    db := GetDB()
    n.Read = true
    db.C("notification").Update(bson.M{"id": n.Id}, bson.M{"$set": bson.M{"read": true}})
}

// Receives the topic id and post content, and create all the notifications for each mention/user inserting them on DB
func CreateMentionsNotificationsFromPost(tid, senderSlug, postContent string){
    mentionRegex := regexp.MustCompile("\\B\\@[\\w\\-]+") // @slug....
    userSlugMentions := mentionRegex.FindAllStringSubmatch(postContent, -1) // array with user slugs

    now := time.Now().Unix() * 1000

    for _, uslug := range userSlugMentions{
        mentionedUser := User{Username: uslug[0][1:]}
        mentionedUser.GenerateSlug() // To generate the correct slug
        mentionedUser, err := GetUserBySlug(mentionedUser.Slug) // removing first "@"
        if err != nil{
            continue
        }
        if mentionedUser.Slug == senderSlug{
            // avoid "self-notifications"
            continue
        }

        n := Notification{
            Type: "mention",
            Params: []string{tid, senderSlug},
            Uid: mentionedUser.Id,
            Read: false,
            Creation_Date: now,
        }
        AddNotification(n)
    }
}

func (n *Notification) GetAllEntities() (map[string]interface{}){
    // get all entities for the notification (topic, users...)
    // depending on the type
    if n.Type == "mention"{
        return n.GetAllEntitiesForMention()
    }else{
        return make(map[string]interface{}, 0)
    }
}

func (n *Notification) GetAllEntitiesForMention() (map[string]interface{}){

    notifJson := make(map[string]interface{}, 0)
    tmp, _ := json.Marshal(n)
    err := json.Unmarshal(tmp, &notifJson)
    if err != nil{
        return notifJson
    }

    entities := make(map[string]interface{}, 0)

    // Get topic info
    entities["topic"], _ = GetTopicById(n.Params[0])
    // Get sender user info
    entities["user"], _ = GetUserBySlugSafe(n.Params[1])

    notifJson["entities"] = entities

    return notifJson

}

