package main

import (
    //"log"
    //"time"
    //"os" 
    //"io"
    //"strings"
    "fmt"
    //"net/http"
    //"html"
    //"encoding/json"
    "flag"

    "GoBBit/server"
    "GoBBit/db"
)

func main(){
    // Command line options
    port := flag.String("p", "", "Set server port")
    staticPath := flag.String("static", "", "Set static resources path (to serve html, js, css.. files (in production you should use nginx or similar))")
    flag.Parse()

	fmt.Printf("Starting..\n")
    
    db.EnsureIndex()
	server.ListenAndServe(*port, *staticPath)

}















