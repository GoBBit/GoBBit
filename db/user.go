package db


import (
	"gopkg.in/mgo.v2/bson"

)


type User struct{
    Id bson.ObjectId `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Picture string `json:"picture"`
    Followed_Communities []bson.ObjectId `json:"followed_communities"`
    Last_Post_Time int64 `json:"last_post_time"`
    Last_Online_Time int64 `json:"last_online_time"`
    Creation_Date int64 `json:"creation_date"`
    Posts_Number int64 `json:"post_number"`
    Topics_Number int64 `json:"topic_number"`
    IsAdmin bool `json:"isadmin"`
}


func AddUser(u User) (User, error){
	db := GetDB()
	
    u.Id = bson.NewObjectId()
    err := db.C("user").Insert(u)

	return u, err
}

func UpdateUserLastPost(id string, t int64) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{ "$set": bson.M{"last_post_time": t} })

    return err
}

func UpdateUserLastOnline(id string, t int64) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{ "$set": bson.M{"last_online_time": t} })

    return err
}

func GetUserByPassword(password string) (User, error){
    db := GetDB()
    
    u := User{}
    err := db.C("user").Find(bson.M{"password":password}).One(&u)

    return u, err
}

func GetUserById(id string) (User, error){
    db := GetDB()
    
    u := User{}
    err := db.C("user").Find(bson.M{"id":bson.ObjectIdHex(id)}).One(&u)

    return u, err
}

func GetUserByIdSafe(id string) (User, error){
    db := GetDB()
    
    u := User{}
    err := db.C("user").Find(bson.M{"id":bson.ObjectIdHex(id)}).Select(bson.M{"id":1, "username": 1, "fid":1, "picture":1}).One(&u)

    return u, err
}

func GetUsersByIds(ids []bson.ObjectId) ([]User, error){
    db := GetDB()
    
    u := []User{}
    err := db.C("user").Find(bson.M{"id": bson.M{"$in": ids }}).All(&u)

    return u, err
}

func GetUsersByIdsSafe(ids []bson.ObjectId) ([]User, error){
    db := GetDB()
    
    u := []User{}
    err := db.C("user").Find(bson.M{"id": bson.M{"$in": ids }}).Select(bson.M{"id":1, "username": 1, "fid":1, "picture":1}).All(&u)

    return u, err
}

func GetUserBySession(id string) (User, error){
    db := GetDB()
    
    u := User{}
    s := UserSession{}
    err := db.C("session").Find(bson.M{"id":id}).One(&s)
    if err != nil{
        return u, err
    }

    err2 := db.C("user").Find(bson.M{"id":s.Uid}).One(&u)

    return u, err2
}

func AddFollowedCommunityToUser(id, cid string) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$push": bson.M{"followed_communities": bson.ObjectIdHex(cid)}})

    return err
}

func DeleteFollowedCommunityToUser(id, cid string) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$pull": bson.M{"followed_communities": bson.ObjectIdHex(cid)}})

    return err
}
















