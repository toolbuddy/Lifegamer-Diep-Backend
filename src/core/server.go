package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"github.com/f26401004/Lifegamer-Diep-backend/src/game"
)

// define the serverHandler struct
type serverHandler struct {
	app *App
	handler func(*App, http.ResponseWriter, *http.Request)
}

// define the ServeHTTP function in serverHandler
func (sh serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sh.handler(sh.app, w, r)
}
// define the gameWebsocketHandler in core package
func gameWebsocketHandler(app *App, w http.ResponseWriter, r *http.Request) {
	// get the websocket instance
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	// return error message if the websocket handshake not established 
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	// just return if other error happend
	} else if err != nil {
		return
	}

	// get the query
	queries, _ := url.ParseQuery(r.URL.RawQuery)
	if (queries["name"][0] == "") {
		log.Fatal("Invalid player name!")
		return
	}
	// generate the player instance
	player := game.NewPlayer(queries["name"][0])
	// default select the first game instance
	select_game := (*app).Games[0]
	// generate the player session
	session := game.NewSession(ws, player, select_game)
	// add the player session to the game
	select_game.JoinPlayer(session)

	log.Printf("Player %s joined to game", queries["name"][0])
}

// define the staticHandler function in core package
func staticHandler (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/dist/" + r.URL.Path[1:])
}

// define runServer in App
func (app App) runServer () {
	// handle the static file in url "/"
	http.HandleFunc("/", staticHandler)
	// handle the game websocket messaging in url "/game_ws"
	http.Handle("/game_ws", serverHandler { &app, gameWebsocketHandler })
	log.Println("run server on port", (*app.Configuration).Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", (*app.Configuration).Server.Host, (*app.Configuration).Server.Port), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}