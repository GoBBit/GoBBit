package server

import (
	"fmt"
	"net/http"
	"encoding/json"

	"GoBBit/db"
)

type UserUpdate struct{
    Email string
    Picture string
    Password string
    IsAdmin bool
}
func GetMeHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
	if e != nil{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Error: No User")
		return
	}

    if r.Method == "GET"{
    	w.Header().Add("Content-Type", "application/json")
    	json.NewEncoder(w).Encode(user)
        return
    }else if r.Method == "PUT"{
        userUpdate := UserUpdate{}
        err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        if userUpdate.Email != ""{
            user.Email = userUpdate.Email
        }

        if userUpdate.Picture != ""{
            user.Picture = userUpdate.Picture
        }

        db.UpdateUser(user)
        json.NewEncoder(w).Encode(user)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

type UserPasswordUpdate struct{
    NewPassword string
    OldPassword string
}
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
    if e != nil{
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "Error: No User")
        return
    }

    if r.Method == "PUT"{
        userUpdate := UserPasswordUpdate{}
        err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        if !user.ChangePassword(userUpdate.OldPassword, userUpdate.NewPassword){
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "invalid_data")
            return
        }

        db.UpdateUser(user)
        json.NewEncoder(w).Encode(user)
        return
    }else{
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
    }

}

func UserHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){

	slug := r.URL.Query().Get("u") // user slug
	u := db.User{}

	if r.Method == "GET"{
		u, err := db.GetUserBySlugSafe(slug)
		if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }
        json.NewEncoder(w).Encode(u)
        return
	}

	if e != nil{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "error_unauthorized")
		return
	}

	// add new user
	if r.Method == "POST"{
	}

	// Edit/Delete an user
	// Only admins!
	if !user.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
	}

	if r.Method == "PUT"{
		userUpdate := UserUpdate{}
		err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        u, err = db.GetUserBySlug(slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        if userUpdate.Password != ""{
        	u.GeneratePasswordHash(userUpdate.Password)
        }

        if userUpdate.Email != ""{
        	u.Email = userUpdate.Email
        }

        if userUpdate.Picture != ""{
        	u.Picture = userUpdate.Picture
        }

        u.IsAdmin = userUpdate.IsAdmin
        db.UpdateUser(u)
	}else if r.Method == "DELETE"{
        
	}else{
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
	}

	json.NewEncoder(w).Encode(u)

}

type UserBan struct{
    Ban bool
}
func UserBanHandler(w http.ResponseWriter, r *http.Request, user db.User, e error){
	// Ban/Unban user from forum
	slug := r.URL.Query().Get("u") // user slug
	u := db.User{}

	// Only admins!
	if !user.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "error_unauthorized")
        return
	}

	if r.Method == "PUT"{
		userUpdate := UserBan{}
		err := json.NewDecoder(r.Body).Decode(&userUpdate)
        if err != nil{
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "invalid_data")
            return
        }
        u, err = db.GetUserBySlug(slug)
        if err != nil{
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "error_user_not_found")
            return
        }

        u.IsBanned = userUpdate.Ban
        db.UpdateUser(u)
	}else{
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: Wrong Method")
        return
	}

	json.NewEncoder(w).Encode(u)

}
