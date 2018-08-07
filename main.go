package main

import (
	"log"
	"net/http"

	"github.com/domino14/aerobot/botui"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

func main() {
	http.HandleFunc("/botui", botui.IndexHandler)
	http.Handle("/", http.FileServer(http.Dir("./botui/static")))

	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(new(botui.AerobotService), "AerobotService")
	http.Handle("/rpc", s)
	log.Fatal(http.ListenAndServe(":8081", nil))
	log.Println("ok")
}

/**
{"jsonrpc":"2.0","method":"AerobotService.Start","id":1,"params":{"username":"","password":"","channel":"","lexiconDb":""}}

*/
