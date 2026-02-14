package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer con.Close()

	for {
		_, msg, err := con.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Received: %s", msg)

		err = con.WriteMessage(websocket.TextMessage, []byte("Hello, client!"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
