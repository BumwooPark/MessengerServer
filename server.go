package main

import (
	_ "encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := flag.String("addr", ":44444", "http service address")
	nsqAddr := flag.String("nsqaddr","54.95.36.34:32771","Message Queue Ipaddress With Port ex)127.0.0.1:30000")
	flag.Parse()

	router := mux.NewRouter()

	//router.HandleFunc("/{userID}/{room:[0-9a-zA-Z]+}",func (w http.ResponseWriter, r *http.Request){
	router.HandleFunc("/{userID}/{room}", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Error", r)
				w.Write([]byte("error"))
			}
		}()

		vars := mux.Vars(r)
		room := vars["room"]
		id := vars["userID"]

		// 채팅방 및 유저 아이디 검사
		userID, roomID, err := db.searchIDAndRoom(id, room)
		if err != nil {
			panic(err)
		}

		serveWs(userID, roomID, *nsqAddr, w, r)
	})

	http.Handle("/", router)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
