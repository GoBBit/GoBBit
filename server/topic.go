package server

import (
    "time"
	"fmt"
	"net/http"
	"encoding/json"
    "strconv"
    "html"

	"GoBBit/db"
    "GoBBit/config"
)

type TopicCreation struct{
    Title string
    Content string // main post content
    Community string // community slug
}
func TopicHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    tid := r.URL.Query().Get("tid") // topic id
    topic := db.Topic{}

    if r.Method == "GET"{
        topic, err := db.GetTopicById(tid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_topic_not_found")
            return
        }

        // now lets add the user creator info to the topic
        tmp, _ := json.Marshal(topic)
        myJson := make(map[string]interface{}, 0)
        _ = json.Unmarshal(tmp, &myJson)
        myJson["user"], _ = db.GetUserByIdSafe(topic.Uid.Hex())

        json.NewEncoder(w).Encode(myJson)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "POST"{
        topicUpdate := TopicCreation{}
        err := json.NewDecoder(r.Body).Decode(&topicUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        // Check if user can post in this community
        community, err := db.GetCommunityBySlug(topicUpdate.Community)
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
        if topicUpdate.Title == "" || topicUpdate.Title == " " || len(topicUpdate.Title) > config.GetInstance().MaxTitleLength || len(topicUpdate.Title) < config.GetInstance().MinTitleLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_title")
            return
        }
        if topicUpdate.Content == "" || topicUpdate.Content == " " || len(topicUpdate.Content) > config.GetInstance().MaxContentLength || len(topicUpdate.Content) < config.GetInstance().MinContentLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_content")
            return
        }

        topic.Title = html.EscapeString(topicUpdate.Title)
        topic.Content = html.EscapeString(topicUpdate.Content)
        topic.Community = topicUpdate.Community
        topic.Uid = user.Id
        topic.GenerateSlug()
        
        now := time.Now().Unix() * 1000
        topic.Creation_Date = now
        topic.Editation_Date = now
        topic.Last_Update = now

        topic, err = db.AddTopic(topic)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        // create notifications
        go db.CreateMentionsNotificationsFromPost(topic.Id.Hex(), user.Slug, topic.Content)

        // update user stats
        db.IncrementPostsNumber(user.Id.Hex(), 1)
        db.IncrementTopicsNumber(user.Id.Hex(), 1)
        db.UpdateUserLastPost(user.Id.Hex(), now)

        // update community stats
        community.IncrementTopicsNumber(1)

        json.NewEncoder(w).Encode(topic)
        return
    }
    
    topic, err := db.GetTopicById(tid)
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
    userCanEdit := (user.IsAdmin || topic.IsOwner(user) || community.IsMod(user))
    if !userCanEdit{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_unauthorized")
        return
    }

    if r.Method == "PUT"{
        topicUpdate := TopicCreation{}
        err := json.NewDecoder(r.Body).Decode(&topicUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        // Security checks
        if topicUpdate.Title == "" || topicUpdate.Title == " " || len(topicUpdate.Title) > config.GetInstance().MaxTitleLength || len(topicUpdate.Title) < config.GetInstance().MinTitleLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_title")
            return
        }
        if topicUpdate.Content == "" || topicUpdate.Content == " " || len(topicUpdate.Content) > config.GetInstance().MaxContentLength || len(topicUpdate.Content) < config.GetInstance().MinContentLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_content")
            return
        }

        topic.Title = html.EscapeString(topicUpdate.Title)
        topic.Content = html.EscapeString(topicUpdate.Content)
        topic.GenerateSlug()

        if user.IsAdmin{
            // only admins can change the community of a topic
            // Check if user can post in this community
            _, err := db.GetCommunityBySlug(topicUpdate.Community)
            if err != nil{
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprintf(w, "error_community_not_found")
                return
            }
            topic.Community = topicUpdate.Community
        }

        now := time.Now().Unix() * 1000
        topic.Editation_Date = now

        db.UpdateTopic(topic)
        json.NewEncoder(w).Encode(topic)
        return
    }else if r.Method == "DELETE"{
        db.DeleteTopic(tid)
        // update user stats
        db.IncrementPostsNumber(user.Id.Hex(), -1)
        db.IncrementTopicsNumber(user.Id.Hex(), -1)
        community.IncrementTopicsNumber(-1)

        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func TopicPostsHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    tid := r.URL.Query().Get("tid") // topic id
    start, _ := strconv.Atoi(r.URL.Query().Get("start"))

    if r.Method == "GET"{
        posts, err := db.GetPostsByTopicIdWithoutIgnored(tid, config.GetInstance().PostsPerPage, start, user.Ignored_Users)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_posts_not_found")
            return
        }

        // now lets add the user creator info to the topic
        tmp, _ := json.Marshal(posts)
        myJson := make([]map[string]interface{}, 0)
        _ = json.Unmarshal(tmp, &myJson)
        for i, p := range posts{
            myJson[i]["user"], _ = db.GetUserByIdSafe(p.Uid.Hex())
        }

        json.NewEncoder(w).Encode(myJson)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "POST"{
    }else if r.Method == "PUT"{
    }else if r.Method == "DELETE"{
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}



func TopicsRecentHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    // Recent topics
    
    start, _ := strconv.Atoi(r.URL.Query().Get("start")) // get from topic num

    if r.Method == "GET"{
        topics, err := db.GetTopicsListWithoutIgnoredUsers(config.GetInstance().TopicsPerPage, start, user.Ignored_Users)
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




