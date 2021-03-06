package server

import (
    "time"
	"fmt"
	"net/http"
	"encoding/json"
    "strconv"
    "html"

    "github.com/tv42/slug"

	"GoBBit/db"
    "GoBBit/config"
)

type CommunityCreation struct{
    Name string
    Description string
    Picture string // picture url
}
func CommunityHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    communitySlug := r.URL.Query().Get("c") // community slug
    community := db.Community{}

    if r.Method == "GET"{
        community, err := db.GetCommunityBySlug(communitySlug)
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

        // Security checks
        if communityUpdate.Name == "" || communityUpdate.Name == " " || len(communityUpdate.Name) > config.GetInstance().MaxNameLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_name")
            return
        }
        if communityUpdate.Description == "" || communityUpdate.Description == " " || len(communityUpdate.Description) > config.GetInstance().MaxDescriptionLength || len(communityUpdate.Description) < config.GetInstance().MinDescriptionLength{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "error_invalid_content")
            return
        }

        community.Name = html.EscapeString(communityUpdate.Name)
        community.Description = html.EscapeString(communityUpdate.Description)
        community.Picture = html.EscapeString(communityUpdate.Picture)
        
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
    community, err := db.GetCommunityBySlug(communitySlug)
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
        community, err = db.GetCommunityBySlug(communitySlug)
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
        db.DeleteCommunityBySlug(communitySlug)
        fmt.Fprintf(w, "ok")
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func CommunityModsHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    communitySlug := r.URL.Query().Get("c") // community slug
    community, err := db.GetCommunityBySlug(communitySlug)
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

    modUserSlug := slug.Slug(r.URL.Query().Get("u")) // user slug to add or remove
    u, err := db.GetUserBySlug(modUserSlug)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_user_not_found")
        return
    }
    modUid := u.Id.Hex()

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

    communitySlug := r.URL.Query().Get("c") // community slug
    community, err := db.GetCommunityBySlug(communitySlug)
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

    banUserSlug := slug.Slug(r.URL.Query().Get("u")) // user slug to add or remove
    u, err := db.GetUserBySlug(banUserSlug)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_user_not_found")
        return
    }
    banUid := u.Id.Hex()

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


func CommunityTopicsHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    communitySlug := r.URL.Query().Get("c") // community slug
    start, _ := strconv.Atoi(r.URL.Query().Get("start")) // get from topic num
    _, err := db.GetCommunityBySlug(communitySlug)
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "error_community_not_found")
        return
    }

    if r.Method == "GET"{
        topics, err := db.GetTopicsByCommunityWithoutIgnoredUsers([]string{communitySlug}, config.GetInstance().TopicsPerPage, start, user.Ignored_Users)
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

func CommunitiesHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    communities := []db.Community{}

    if r.Method == "GET"{
        communities, _ = db.GetAllCommunities()
        json.NewEncoder(w).Encode(communities)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }
}





