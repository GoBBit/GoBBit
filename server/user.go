package server

import (
	"fmt"
	"net/http"
	"encoding/json"
    "strconv"
    "html"

	"GoBBit/db"
)

type UserUpdate struct{
    Email string
    Picture string
    Password string
    IsAdmin bool
}
func GetMeHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
	if e != nil{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Error: No User")
		return
	}

    if r.Method == "GET"{
    	w.Header().Add("Content-Type", "application/json")
    	json.NewEncoder(w).Encode(user)
        return
    }else if r.Method == "PUT"{
        userUpdate := UserUpdate{}
        err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        if userUpdate.Email != ""{
            user.Email = html.EscapeString(userUpdate.Email)
        }

        if userUpdate.Picture != ""{
            user.Picture = html.EscapeString(userUpdate.Picture)
        }

        db.UpdateUser(user)
        json.NewEncoder(w).Encode(user)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

type UserPasswordUpdate struct{
    NewPassword string
    OldPassword string
}
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "PUT"{
        userUpdate := UserPasswordUpdate{}
        err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        if !user.ChangePassword(userUpdate.OldPassword, userUpdate.NewPassword){
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        db.UpdateUser(user)
        json.NewEncoder(w).Encode(user)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func UserHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

	slug := r.URL.Query().Get("u") // user slug
	u := db.User{}
    var err error

	if r.Method == "GET"{
        u.Username = slug
        u.GenerateSlug()
        if user.IsAdmin || user.Slug == slug{
            u, err = db.GetUserBySlug(u.Slug)
        }else{
    		u, err = db.GetUserBySlugSafe(u.Slug)
        }

        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        json.NewEncoder(w).Encode(u)
        return
	}

	if e != nil{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "error_unauthorized")
		return
	}

	// add new user
	if r.Method == "POST"{
	}

	// Edit/Delete an user
	// Only admins!
	if !user.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
	}

	if r.Method == "PUT"{
		userUpdate := UserUpdate{}
		err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        u, err = db.GetUserBySlug(slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        if userUpdate.Password != ""{
        	u.GeneratePasswordHash(userUpdate.Password)
        }

        if userUpdate.Email != ""{
        	u.Email = userUpdate.Email
        }

        if userUpdate.Picture != ""{
        	u.Picture = userUpdate.Picture
        }

        u.IsAdmin = userUpdate.IsAdmin
        db.UpdateUser(u)
	}else if r.Method == "DELETE"{
        
	}else{
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
	}

	json.NewEncoder(w).Encode(u)

}

type UserBan struct{
    Ban bool
}
func UserBanHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
	// Ban/Unban user from forum
	uid := r.URL.Query().Get("uid") // user slug

	// Only admins!
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }
	if !user.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
	}

    u, err := db.GetUserById(uid)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_user_not_found")
        return
    }

    if r.Method == "POST"{
        u.IsBanned = true
        db.UpdateUser(u)
        
        json.NewEncoder(w).Encode(user)
        return
    }else if r.Method == "DELETE"{
        u.IsBanned = false
        db.UpdateUser(u)
        
        json.NewEncoder(w).Encode(user)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}


func UserHomeHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    // User Home: topics based on the communities user follows
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    start, _ := strconv.Atoi(r.URL.Query().Get("start")) // get from topic num

    if r.Method == "GET"{
        topics, err := db.GetTopicsByCommunityWithoutIgnoredUsers(user.Followed_Communities, TopicsPerPage, start, user.Ignored_Users)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_topics_not_found")
            return
        }

        // now lets add the user creator info to the topic
        tmp, _ := json.Marshal(topics)
        myJson := make([]map[string]interface{}, 0)
        _ = json.Unmarshal(tmp, &myJson)
        for i, t := range topics{
            myJson[i]["user"], _ = db.GetUserByIdSafe(t.Uid.Hex())
        }

        json.NewEncoder(w).Encode(myJson)
        return

    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func UserFollowCommunityHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    // User Home: topics based on the communities user follows
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    cslug := r.URL.Query().Get("c") // community slug to delete from followed communities

    if r.Method == "POST"{
        community, err := db.GetCommunityBySlug(cslug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_community_not_found")
            return
        }

        db.DeleteFollowedCommunityToUser(user.Id.Hex(), community.Slug) // delete if already followed
        err = db.AddFollowedCommunityToUser(user.Id.Hex(), community.Slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_community_not_followed")
            return
        }

        json.NewEncoder(w).Encode(community)
        return

    }else if r.Method == "DELETE"{
        db.DeleteFollowedCommunityToUser(user.Id.Hex(), cslug)

        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func IgnoreUserHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "GET"{
        w.Header().Add("Content-Type", "application/json")
        ignored_users, _ := db.GetUsersByIdsSafe(user.Ignored_Users)
        json.NewEncoder(w).Encode(ignored_users)
        return
    }else if r.Method == "POST"{
        uid := r.URL.Query().Get("uid") // user id
        db.DeleteIgnoredUserToUser(user.Id.Hex(), uid) // delete if already exists
        db.AddIgnoredUserToUser(user.Id.Hex(), uid)
        
        json.NewEncoder(w).Encode(user)
        return
    }else if r.Method == "DELETE"{
        uid := r.URL.Query().Get("uid") // user id
        db.DeleteIgnoredUserToUser(user.Id.Hex(), uid)
        
        json.NewEncoder(w).Encode(user)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}


func UserTopicsHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    slug := r.URL.Query().Get("u") // user slug
    start, _ := strconv.Atoi(r.URL.Query().Get("start")) // get from topic num
    u := db.User{Username: slug}
    u.GenerateSlug()

    if r.Method == "GET"{
        u, err := db.GetUserBySlugSafe(u.Slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        topics, err := db.GetTopicsByUser(u.Id.Hex(), TopicsPerPage, start)

        json.NewEncoder(w).Encode(topics)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

