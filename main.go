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
    flag.Parse()

	fmt.Printf("Starting..\n")
    
    db.EnsureIndex()
	server.ListenAndServe(*port)

}















