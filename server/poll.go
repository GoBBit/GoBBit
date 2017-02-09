package server

import (
    "fmt"
    "time"
	"net/http"
	"encoding/json"
    "html"

	"GoBBit/db"
    "GoBBit/config"
)

type PollCreation struct{
    Title string `json:"title"`
    Options []string `json:"options"`
    Tid string `json:"tid"`
    Allowed_Options int64
}

func PollHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    pollid := r.URL.Query().Get("pollid")

    if r.Method == "GET"{
        poll, err := db.GetPollById(pollid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_poll_not_found")
            return
        }

        json.NewEncoder(w).Encode(poll)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "POST"{
        newPoll := PollCreation{}
        err := json.NewDecoder(r.Body).Decode(&newPoll)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        // Check permissions
        topic, err := db.GetTopicById(newPoll.Tid)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
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

        // Security checks
        if newPoll.Title == "" || newPoll.Title == " " || len(newPoll.Title) > config.GetInstance().MaxTitleLength || len(newPoll.Title) < config.GetInstance().MinTitleLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_title")
            return
        }

        poll := db.Poll{}
        poll.Title = html.EscapeString(newPoll.Title)
        poll.Tid = topic.Id
        poll.Allowed_Options = 1 //newPoll.Allowed_Options

        now := time.Now().Unix() * 1000
        poll.Creation_Date = now

        for _,o := range newPoll.Options{
            pollOption := db.PollOption{Title: html.EscapeString(o)}
            poll.Options = append(poll.Options, pollOption)
        }

        poll, err = db.AddPoll(poll)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }

    	json.NewEncoder(w).Encode(poll)
        return
    }else if r.Method == "DELETE"{
        // Check permissions
        topic, err := db.GetTopicById(pollid)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
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

        db.DeletePoll(pollid)
        fmt.Fprintf(w, "ok")
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }
}

type Vote struct{
    Title string `json:"title"`
    PollId string
}

func PollVoteHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    vote := Vote{}
    err := json.NewDecoder(r.Body).Decode(&vote)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "invalid_data")
        return
    }

    if r.Method == "POST"{
        poll, err := db.GetPollById(vote.PollId)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_poll_not_found")
            return
        }

        for _, o := range poll.Options{
            // clean older votes if exists
            db.DeleteVoteToPoll(vote.PollId, o.Title, user.Id.Hex())
        }
        db.AddVoteToPoll(vote.PollId, vote.Title, user.Id.Hex())

        json.NewEncoder(w).Encode(poll)
        return
    }else if r.Method == "DELETE"{
        db.DeleteVoteToPoll(vote.PollId, vote.Title, user.Id.Hex())
        fmt.Fprintf(w, "ok")
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }
}



