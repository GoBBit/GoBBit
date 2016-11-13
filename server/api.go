package server

import (
	"log"
	"time"
	"os" 
	//"io"
	//"io/ioutil"
	//"strings"
	"fmt"
	"net/http"
	//"net/url"
	//"html"
	"encoding/json"

	"GoBBit/db"
	"GoBBit/utils"
)

var (
	Port = ":3000"
)

func ConfigEnvVars(){

	if os.Getenv("PORT") != ""{
		Port = ":" + os.Getenv("PORT")
	}
}

func ListenAndServe(cmdPort string){
	ConfigEnvVars()
	if cmdPort != ""{
		Port = ":" + cmdPort
	}

	// Setup routes
	mux := http.NewServeMux()

    // User Endpoints
    mux.HandleFunc("/api/me", Middleware(GetMeHandler))
    mux.HandleFunc("/api/user", Middleware(UserHandler))
    mux.HandleFunc("/api/user/changepassword", Middleware(ChangePasswordHandler))
	mux.HandleFunc("/api/user/ban", Middleware(UserBanHandler))

    // Topic Endpoints
    mux.HandleFunc("/api/topic", Middleware(TopicHandler))

    // Community Endpoints
    mux.HandleFunc("/api/community", Middleware(CommunityHandler))
    mux.HandleFunc("/api/community/mods", Middleware(CommunityModsHandler))
    mux.HandleFunc("/api/community/ban", Middleware(CommunityBannedUsersHandler))

	// Login & LogOut
	mux.HandleFunc("/register", Middleware(RegisterHandler))
    mux.HandleFunc("/login", Middleware(LoginHandler))
	mux.HandleFunc("/logout", Middleware(LogoutHandler))
	
	// mux.Handle("/", http.FileServer(http.Dir("./public_html")))
	mux.Handle("/debug/vars", http.DefaultServeMux)

	fmt.Printf("listening on *%s\n", Port)

	// Start listening
	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		Addr:         Port,
		Handler:      mux,
	}
	log.Fatal(s.ListenAndServe())
}


func Middleware(next func(http.ResponseWriter, *http.Request, db.User, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("xsession")
		if err != nil{
			next(w, r, db.User{}, err)
			return
		}

		u, err := db.GetUserBySession(cookie.Value)

		next(w, r, u, err)
	}
}

type RegisterUser struct{
    Username string
    Password string
    Email string
}
func RegisterHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if r.Method != "POST"{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: Wrong method")
        return
    }

    rUser := RegisterUser{}
    err := json.NewDecoder(r.Body).Decode(&rUser)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Unable to decode user")
        return
    }

    u := db.User{}
    u.Username = rUser.Username
    u.GeneratePasswordHash(rUser.Password)
    u.Email = rUser.Email
    u.GenerateSlug()

    now := time.Now().Unix() * 1000
    u.Creation_Date = now
    u.Last_Post_Time = now
    u.Last_Online_Time = now

    u, err2 := db.AddUser(u)
    if err2 != nil{
        fmt.Fprintf(w, "Error: Unable to create user", err2)
        return
    }

    json.NewEncoder(w).Encode(u)

}

type LoginUser struct{
    Username string
    Password string
}
func LoginHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if r.Method != "POST"{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: Wrong method")
        return
    }

    rUser := LoginUser{}
    err := json.NewDecoder(r.Body).Decode(&rUser)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Unable to decode user")
        return
    }

    u := db.User{Username:rUser.Username}
    hash := u.GeneratePasswordHash(rUser.Password)
    u, err2 := db.GetUserByPassword(hash)
    if err2 != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Error: Unable to get user")
        return
    }

    // Create session
	sessionHash := utils.GenerateUserSession()
	uSess := db.UserSession{Uid:u.Id, Id:sessionHash}
    us, err3 := db.AddUserSession(uSess)
    if err3 != nil{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Unable to create session [%s]", err3)
        return
    }

    xsess := "xsession=" + us.Id
    cookie := http.Cookie{Name:"xsession",Value:us.Id, Path:"/",Expires:time.Now().AddDate(1,0,0)}

    http.SetCookie(w, &cookie)
    w.Header().Add("Cookie", xsess)
    json.NewEncoder(w).Encode(u)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    
    actCookie, err := r.Cookie("xsession")
    if err != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Error: No cookie found")
        return
    }

    db.DeleteUserSession(actCookie.Value)

    xsess := "xsession=;"
    cookie := http.Cookie{Name:"xsession",Value:"", Expires:time.Now().AddDate(0,0,-1)}

    http.SetCookie(w, &cookie)
    w.Header().Add("Cookie", xsess)
    
    fmt.Fprintf(w, "ok")

}


