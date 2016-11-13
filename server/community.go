package server

import (
    "time"
	"fmt"
	"net/http"
	"encoding/json"

	"GoBBit/db"
)

type CommunityCreation struct{
    Name string
    Description string
    Picture string // picture url
}
func CommunityHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    slug := r.URL.Query().Get("c") // community slug
    community := db.Community{}

    if r.Method == "GET"{
        community, err := db.GetCommunityBySlug(slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_community_not_found")
            return
        }

    	json.NewEncoder(w).Encode(community)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "POST"{
        // Only admins can create new communities
        if !user.IsAdmin{
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "error_unauthorized")
            return
        }

        communityUpdate := CommunityCreation{}
        err := json.NewDecoder(r.Body).Decode(&communityUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        community.Name = communityUpdate.Name
        community.Description = communityUpdate.Description
        community.Picture = communityUpdate.Picture
        
        now := time.Now().Unix() * 1000
        community.Creation_Date = now
        community.GenerateSlug()

        community, err = db.AddCommunity(community)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        json.NewEncoder(w).Encode(community)
        return
    }

    // Only admins or mods can update/delete communities 
    community, err := db.GetCommunityBySlug(slug)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_community_not_found")
        return
    }
    if !user.IsAdmin && !community.IsMod(user){
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
    }

    if r.Method == "PUT"{
        communityUpdate := CommunityCreation{}
        err := json.NewDecoder(r.Body).Decode(&communityUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        community, err = db.GetCommunityBySlug(slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_community_not_found")
            return
        }

        community.Description = communityUpdate.Description
        community.Picture = communityUpdate.Picture

        err = db.UpdateCommunity(community)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        json.NewEncoder(w).Encode(community)
        return
    }else if r.Method == "DELETE"{
        db.DeleteCommunityBySlug(slug)
        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func CommunityModsHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    slug := r.URL.Query().Get("c") // community slug
    community, err := db.GetCommunityBySlug(slug)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_community_not_found")
        return
    }

    if r.Method == "GET"{
        mods, _ := db.GetUsersByIdsSafe(community.Mods)
        json.NewEncoder(w).Encode(mods)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    // Only admins or mods can add/delete mods
    if !user.IsAdmin && !community.IsMod(user){
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
    }

    modUid := r.URL.Query().Get("uid") // mod uid to add or remove
    if r.Method == "POST"{
        db.DeleteModsToCommunity(community.Id.Hex(), modUid) // delete mod if already added
        err := db.AddModsToCommunity(community.Id.Hex(), modUid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "invalid_data %s", err)
            return
        }
        json.NewEncoder(w).Encode(community)
        return
    }else if r.Method == "DELETE"{
        db.DeleteModsToCommunity(community.Id.Hex(), modUid)
        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func CommunityBannedUsersHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    slug := r.URL.Query().Get("c") // community slug
    community, err := db.GetCommunityBySlug(slug)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_community_not_found")
        return
    }

    if r.Method == "GET"{
        mods, _ := db.GetUsersByIdsSafe(community.Banned_Users)
        json.NewEncoder(w).Encode(mods)
        return
    }

    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    // Only admins or mods can add/delete mods
    if !user.IsAdmin && !community.IsMod(user){
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
    }

    banUid := r.URL.Query().Get("uid") // user uid to add or remove
    if r.Method == "POST"{
        db.DeleteBannedUserToCommunity(community.Id.Hex(), banUid) // delete if already added
        err := db.AddBannedUserToCommunity(community.Id.Hex(), banUid)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "invalid_data %s", err)
            return
        }
        json.NewEncoder(w).Encode(community)
        return
    }else if r.Method == "DELETE"{
        db.DeleteBannedUserToCommunity(community.Id.Hex(), banUid)
        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}







