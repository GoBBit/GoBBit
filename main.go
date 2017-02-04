package main

import (
    "fmt"
    "flag"

    "GoBBit/server"
    "GoBBit/db"
    "GoBBit/config"
)

// make admin by userslug on command line
// it is useful to create the first admin in the forum
func makeAdmin(userslug string){
    user := db.User{Username: userslug}
    user.GenerateSlug()
    user.MakeAdmin(true)
}

func main(){
    // Command line options
    port := flag.String("p", "", "Set server port")
    staticPath := flag.String("static", "", "Set static resources path (to serve html, js, css.. files (in production you should use nginx or similar))")
    cfgFile := flag.String("c", "./config.json", "Config file")
    mAdmin := flag.String("admin", "", "Make admin an user by userslug")
    flag.Parse()

    fmt.Printf("\nLoading config.json..")
    config.CreateInstance(*cfgFile)
    if *port != ""{
        config.GetInstance().Port = *port
    }
    
    // Prepare DB
    db.EnsureIndex()

    if *mAdmin != ""{
        makeAdmin(*mAdmin)
        return
    }

	server.ListenAndServe(*staticPath)

}

