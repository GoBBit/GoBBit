package server

import (
    "time"
	"fmt"
	"net/http"
	"encoding/json"
    "html"

	"GoBBit/db"
)

type PostCreation struct{
    Tid string // topic id
    Content string
}
func PostHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    pid := r.URL.Query().Get("pid") // topic id
    post := db.Post{}

    if r.Method == "GET"{
        post, err := db.GetPostById(pid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_post_not_found")
            return
        }

    	json.NewEncoder(w).Encode(post)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "POST"{
        postUpdate := PostCreation{}
        err := json.NewDecoder(r.Body).Decode(&postUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        // Check if user can post in this community
        topic, err := db.GetTopicById(postUpdate.Tid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_topic_not_found")
            return
        }
        community, err := db.GetCommunityBySlug(topic.Community)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_community_not_found")
            return
        }
        if !community.UserCanPost(user){
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_unauthorized")
            return
        }

        // Security checks
        if postUpdate.Content == "" || postUpdate.Content == " " || len(postUpdate.Content) > MaxContentLength || len(postUpdate.Content) < MinContentLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_content")
            return
        }

        post.Content = html.EscapeString(postUpdate.Content)
        post.Tid = topic.Id
        post.Uid = user.Id
        
        now := time.Now().Unix() * 1000
        post.Creation_Date = now
        post.Editation_Date = now

        post, err = db.AddPost(post)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data %s", err)
            return
        }
        db.UpdateTopicLastUpdate(topic.Id.Hex(), now)

        // update user stats
        db.IncrementPostsNumber(user.Id.Hex(), 1)
        db.UpdateUserLastPost(user.Id.Hex(), now)

        // update topic stats
        topic.IncrementPostsNumber(1)

        // update community stats
        community.IncrementPostsNumber(1)

        json.NewEncoder(w).Encode(post)
        return
    }
    
    post, err := db.GetPostById(pid)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_post_not_found")
        return
    }
    topic, err := db.GetTopicById(post.Tid.Hex())
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_topic_not_found")
        return
    }
    community, err := db.GetCommunityBySlug(topic.Community)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_topic_not_found")
        return
    }
    userCanEdit := (user.IsAdmin || post.IsOwner(user) || community.IsMod(user))
    if !userCanEdit{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_unauthorized")
        return
    }

    if r.Method == "PUT"{
        postUpdate := PostCreation{}
        err := json.NewDecoder(r.Body).Decode(&postUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        // Security checks
        if postUpdate.Content == "" || postUpdate.Content == " " || len(postUpdate.Content) > MaxContentLength || len(postUpdate.Content) < MinContentLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_content")
            return
        }

        post.Content = html.EscapeString(postUpdate.Content)
        now := time.Now().Unix() * 1000
        post.Editation_Date = now

        db.UpdatePost(post)
        json.NewEncoder(w).Encode(post)
        return
    }else if r.Method == "DELETE"{
        db.DeletePost(pid)
        // update user stats
        db.IncrementPostsNumber(user.Id.Hex(), -1)
        topic.IncrementPostsNumber(-1)
        community.IncrementPostsNumber(-1)

        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}








