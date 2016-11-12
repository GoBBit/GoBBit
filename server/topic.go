package server

import (
    "time"
	"fmt"
	"net/http"
	"encoding/json"

	"GoBBit/db"
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

    	json.NewEncoder(w).Encode(topic)
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

        topic.Title = topicUpdate.Title
        topic.Content = topicUpdate.Content
        topic.Community = topicUpdate.Community // TODO: Check if community exists
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
        json.NewEncoder(w).Encode(topic)
        return
    }else if r.Method == "PUT"{
        topicUpdate := TopicCreation{}
        err := json.NewDecoder(r.Body).Decode(&topicUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        topic, err := db.GetTopicById(tid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_topic_not_found")
            return
        }

        topic.Title = topicUpdate.Title
        topic.Content = topicUpdate.Content
        topic.Community = topicUpdate.Community // TODO: Check if community exists
        topic.GenerateSlug()

        now := time.Now().Unix() * 1000
        topic.Editation_Date = now

        db.UpdateTopic(topic)
        json.NewEncoder(w).Encode(topic)
        return
    }else if r.Method == "DELETE"{
        db.DeleteTopic(tid)
        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}








