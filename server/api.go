package server

import (
	"log"
	"time"
	"os" 
	//"io"
	//"io/ioutil"
	//"strings"
	"strconv"
	"fmt"
	"net/http"
	//"net/url"
	//"html"
	"encoding/json"
	"encoding/base64"
	"crypto/sha512"

	"gopkg.in/mgo.v2/bson"

	"GoBBit/db"
	//"GoBBit/utils"
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

	mux.HandleFunc("/api/me", Middleware(GetMeHandler))

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


func GetMeHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
	if e != nil{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Error: No User")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	//fmt.Fprintf(w, "Hello, You are looking for %s", r.URL.Query().Get("q"))

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
    u.Password = calculateHash(rUser.Username + rUser.Password)
    u.Email = rUser.Email

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

    hash := calculateHash(rUser.Username + rUser.Password)
    u, err2 := db.GetUserByPassword(hash)
    if err2 != nil{
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Error: Unable to get user")
        return
    }

    // Create session
	sessionHash := GenerateUserSession()
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


// Hash Functions
func calculateHash(s string) (string){
    sha512.New()
    sBytes := []byte(s)
    hash := sha512.Sum512(sBytes)
    b64hash := base64.URLEncoding.EncodeToString(hash[:])

    return b64hash
}

// Session Function
func GenerateUserSession()(string){
	// Generate usersession based on an mongodb objectID and the actual timestamp
	now := time.Now().Unix() * 1000
	nowStr := strconv.FormatInt(now, 10)
	id := bson.NewObjectId()
	session := nowStr + id.Hex()

	return calculateHash(session)
}







