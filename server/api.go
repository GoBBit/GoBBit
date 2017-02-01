package server

import (
	"log"
	"time"
	"os"
    "bytes"
	"io"
	"io/ioutil"
	"strings"
	"fmt"
	"net/http"
    "errors"
	//"net/url"
	"html"
	"encoding/json"

	"GoBBit/db"
	"GoBBit/utils"
)

var (
	Port = ":3000"
    TopicsPerPage = 20
    PostsPerPage = 20
    MaxTitleLength = 140
    MinTitleLength = 5
    MaxContentLength = 1000
    MinContentLength = 5
    MaxNameLength = 50
    MaxDescriptionLength = 140
    MinDescriptionLength = 5
    SITE_KEY = "Change_Me" // SITE_KEY is the key to generate session and other important stuff, please change it on production
)

func ConfigEnvVars(){

	if os.Getenv("PORT") != ""{
		Port = ":" + os.Getenv("PORT")
	}
    if os.Getenv("SITE_KEY") != ""{
        SITE_KEY = os.Getenv("SITE_KEY")
    }
}

func ListenAndServe(cmdPort string, staticPath string){
	ConfigEnvVars()
	if cmdPort != ""{
		Port = ":" + cmdPort
	}

	// Setup routes
	mux := http.NewServeMux()

    // User Endpoints
    mux.HandleFunc("/api/me", Middleware(GetMeHandler))
    mux.HandleFunc("/api/user", Middleware(UserHandler))
    mux.HandleFunc("/api/user/topics", Middleware(UserTopicsHandler))
    mux.HandleFunc("/api/user/changepassword", Middleware(ChangePasswordHandler))
    mux.HandleFunc("/api/user/ban", Middleware(UserBanHandler))
    mux.HandleFunc("/api/user/follow/community", Middleware(UserFollowCommunityHandler))
	mux.HandleFunc("/api/user/home", Middleware(UserHomeHandler))
    mux.HandleFunc("/api/user/ignore", Middleware(IgnoreUserHandler))

    // Notifications Endpoints
    mux.HandleFunc("/api/notifications", Middleware(NotificationsHandler))
    mux.HandleFunc("/api/notification", Middleware(NotificationHandler))
    mux.HandleFunc("/api/notification/read", Middleware(NotificationReadHandler))

    // Post Endpoints
    mux.HandleFunc("/api/post", Middleware(PostHandler))

    // Topic Endpoints
    mux.HandleFunc("/api/topic", Middleware(TopicHandler))
    mux.HandleFunc("/api/topic/posts", Middleware(TopicPostsHandler))
    mux.HandleFunc("/api/topics/recent", Middleware(TopicsRecentHandler))

    // Community Endpoints
    mux.HandleFunc("/api/community", Middleware(CommunityHandler))
    mux.HandleFunc("/api/community/topics", Middleware(CommunityTopicsHandler))
    mux.HandleFunc("/api/community/mods", Middleware(CommunityModsHandler))
    mux.HandleFunc("/api/community/ban", Middleware(CommunityBannedUsersHandler))
    mux.HandleFunc("/api/communities", Middleware(CommunitiesHandler))

	// Login & LogOut
	mux.HandleFunc("/register", Middleware(RegisterHandler))
    mux.HandleFunc("/login", Middleware(LoginHandler))
	mux.HandleFunc("/logout", Middleware(LogoutHandler))
	
    if staticPath != ""{
        // serve static files (html, js...)
        // recommended for development, in production you should use nginx or something like that to serve static files
        mux.Handle("/", http.FileServer(http.Dir(staticPath)))
    }

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


type CSRFRequest struct{
    CSRF string
}
func Middleware(next func(http.ResponseWriter, *http.Request, db.User, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("xsession")
		if err != nil{
			next(w, r, db.User{}, err)
			return
		}


        if r.Method == "POST"{
            // Let's check CSRF token
            // temporary buffer
            b := bytes.NewBuffer(make([]byte, 0))
            // TeeReader returns a Reader that writes to b what it reads from r.Body.
            reader := io.TeeReader(r.Body, b)
            csrf := CSRFRequest{}
            err := json.NewDecoder(reader).Decode(&csrf)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, "error_parsing_csrf")
                return
            }
            if csrf.CSRF != cookie.Value{
                w.WriteHeader(http.StatusUnauthorized)
                fmt.Fprintf(w, "error_bad_csrf")
                return
            }
            // we are done with body
            defer r.Body.Close()
            r.Body = ioutil.NopCloser(b)
        }


        // split by ":" (parse cookie)
        splittedCookie := strings.Split(cookie.Value, ":")
        if len(splittedCookie) < 3{
            next(w, r, db.User{}, errors.New("Invalid Cookie"))
            return
        }

        uid := splittedCookie[0]
        timestamp := splittedCookie[1]
        hash := splittedCookie[2]
        
        u, err := db.GetUserById(uid)
        if err != nil{
            next(w, r, db.User{}, err)
            return
        }
        validSession := utils.CheckSession(uid, u.Password, timestamp, SITE_KEY, hash)
        if validSession{
            next(w, r, u, nil)
            return
        }else{
            next(w, r, db.User{}, errors.New("Invalid Cookie"))
            return
        }
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
    u.Username = html.EscapeString(rUser.Username)
    u.GenerateSlug()
    u.GeneratePasswordHash(rUser.Password)
    u.Email = rUser.Email

    now := time.Now().Unix() * 1000
    u.Creation_Date = now
    u.Last_Post_Time = now
    u.Last_Online_Time = now

    u, err2 := db.AddUser(u)
    if err2 != nil{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Unable to create user")
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

    u := db.User{Username: html.EscapeString(rUser.Username)}
    u.GenerateSlug()
    hash := u.GeneratePasswordHash(rUser.Password)
    u, err2 := db.GetUserByPassword(hash)
    if err2 != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Error: Unable to get user")
        return
    }

    // Create session
	sessionHash := utils.GenerateUserSession(u.Id.Hex(), u.Password, SITE_KEY)

    xsess := "xsession=" + sessionHash
    cookie := http.Cookie{Name:"xsession",Value:sessionHash, Path:"/",Expires:time.Now().AddDate(1,0,0)}

    http.SetCookie(w, &cookie)
    w.Header().Add("Cookie", xsess)
    json.NewEncoder(w).Encode(u)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

    xsess := "xsession=;"
    cookie := http.Cookie{Name:"xsession",Value:"", Expires:time.Now().AddDate(0,0,-1)}

    http.SetCookie(w, &cookie)
    w.Header().Add("Cookie", xsess)
    
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}



