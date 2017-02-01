package server

import (
	"fmt"
	"net/http"
	"encoding/json"

	"GoBBit/db"
)

func NotificationHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    nid := r.URL.Query().Get("nid") // notification id

    if r.Method == "GET"{
        notification, err := db.GetNotificationById(nid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_notification_not_found")
            return
        }

    	json.NewEncoder(w).Encode(notification)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    // list notifications
    if r.Method == "GET"{
        notifications, _ := db.GetNotificationsByUser(user.Id.Hex())
        json.NewEncoder(w).Encode(notifications)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func NotificationReadHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    // mark notification as read
    
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    nid := r.URL.Query().Get("nid") // notification id

    if r.Method == "POST"{
        notification, err := db.GetNotificationById(nid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_notification_not_found")
            return
        }

        if notification.Uid != user.Id{
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "error_unauthorized")
            return
        }

        notification.MarkAsRead()
        json.NewEncoder(w).Encode(notification)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

