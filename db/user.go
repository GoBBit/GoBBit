package db


import (
	"gopkg.in/mgo.v2/bson"
    "github.com/tv42/slug"

    utils "GoBBit/utils"
)


type User struct{
    Id bson.ObjectId `json:"id"`
    Username string `json:"username"`
	Slug string `json:"slug"`
	Email string `json:"email"`
	Password string `json:"password"`
	Picture string `json:"picture"`
    Followed_Communities []string `json:"followed_communities"` // slugs
    Last_Post_Time int64 `json:"last_post_time"`
    Last_Online_Time int64 `json:"last_online_time"`
    Creation_Date int64 `json:"creation_date"`
    Posts_Number int64 `json:"posts_number"`
    Topics_Number int64 `json:"topics_number"`
    IsAdmin bool `json:"isadmin"`
    IsBanned bool `json:"isbanned"`
    Ignored_Users []bson.ObjectId `json:"ignored_users"` // users id array
}


func AddUser(u User) (User, error){
	db := GetDB()
	
    u.Id = bson.NewObjectId()
    err := db.C("user").Insert(u)

	return u, err
}

func UpdateUser(u User) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": u.Id}, u)

    return err
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

func GetUserBySlug(slug string) (User, error){
    db := GetDB()
    
    u := User{}
    err := db.C("user").Find(bson.M{"slug":slug}).One(&u)

    return u, err
}

func GetUserBySlugSafe(slug string) (User, error){
    db := GetDB()
    
    u := User{}
    err := db.C("user").Find(bson.M{"slug":slug}).Select(bson.M{"password":0,"email":0}).One(&u)

    return u, err
}

func GetUserByIdSafe(id string) (User, error){
    db := GetDB()
    
    u := User{}
    err := db.C("user").Find(bson.M{"id":bson.ObjectIdHex(id)}).Select(bson.M{"password":0,"email":0}).One(&u)

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
    err := db.C("user").Find(bson.M{"id": bson.M{"$in": ids }}).Select(bson.M{"password":0,"email":0}).All(&u)

    return u, err
}

func AddFollowedCommunityToUser(id, cslug string) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$push": bson.M{"followed_communities": cslug}})

    return err
}

func DeleteFollowedCommunityToUser(id, cslug string) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$pull": bson.M{"followed_communities": cslug}})

    return err
}

func AddIgnoredUserToUser(id, ignored string) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$push": bson.M{"ignored_users": bson.ObjectIdHex(ignored)}})

    return err
}

func DeleteIgnoredUserToUser(id, ignored string) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$pull": bson.M{"ignored_users": bson.ObjectIdHex(ignored)}})

    return err
}

func IncrementPostsNumber(id string, n int) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$inc": bson.M{"posts_number": n}})

    return err
}

func IncrementTopicsNumber(id string, n int) error{
    db := GetDB()

    err := db.C("user").Update(bson.M{"id": bson.ObjectIdHex(id)}, bson.M{"$inc": bson.M{"topics_number": n}})

    return err
}



func (u *User) GenerateSlug() (string){
    u.Slug = slug.Slug(u.Username)
    return u.Slug
}

func (u *User) GeneratePasswordHash(pass string) (string){
    u.Password = utils.CalculateHash(u.Slug + pass)
    return u.Password
}

func (u *User) ChangePassword(oldpass, newpass string) (bool){
    oldCheck := utils.CalculateHash(u.Slug + oldpass)
    if u.Password != oldCheck{
        return false
    }

    u.Password = utils.CalculateHash(u.Slug + newpass)
    return true
}

func (u *User) MakeAdmin(option bool){
    u.IsAdmin = option

    db := GetDB()
    db.C("user").Update(bson.M{"slug": u.Slug}, bson.M{"$set": bson.M{"isadmin": option}})
}
