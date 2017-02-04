package main

import (
    "fmt"
    "flag"

    "GoBBit/server"
    "GoBBit/db"
    "GoBBit/config"
)

func main(){
    // Command line options
    port := flag.String("p", "", "Set server port")
    staticPath := flag.String("static", "", "Set static resources path (to serve html, js, css.. files (in production you should use nginx or similar))")
    cfgFile := flag.String("c", "./config.json", "Config file")
    flag.Parse()

    fmt.Printf("\nLoading config.json..")
    config.CreateInstance(*cfgFile)
    if *port != ""{
        config.GetInstance().Port = *port
    }

	fmt.Printf("\nStarting..")
    
    db.EnsureIndex()
	server.ListenAndServe(*staticPath)

}

