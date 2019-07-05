package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"net/http/pprof"
	"github.com/f26401004/Lifegamer-Diep-backend/src/game"
)

/**
 * serverHandler:
 * The struct to bind the server handler with app
 *
 * @property {*App} app 																												- the app reference
 * @property {func(*App, http.ResponseWriter, *http.Request)} handler						- the function to handle the join request
 */
type serverHandler struct {
	app *App
	Handler func(*App, http.ResponseWriter, *http.Request)
}

/**
 * <serverHandler>.ServeHTTP:
 * The function in serverHandler to handle the request from client
 *
 * @param {http.ResponseWriter} w 																							- the response writer of current request
 * @param {*http.Request} r																											- the current request
 *
 * @return {nil}
 */
func (sh serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sh.Handler(sh.app, w, r)
}


/**
 * <core>.gameWebsocketHandler:
 * The function to handle the request upgrade to websocket from client
 *
 * @param {*App} app 																														- the app reference
 * @param {http.ResponseWriter} w																								- the response writer of current request
 * @param {*http.Request} r																											- the current request
 *
 * @return {nil}
 */
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
		log.Println("[Error]: Player name can not be empty!")
		http.Error(w, "Player name can not be empty!", 400)
		return
	}
	var select_game *game.Game = nil;
	// search the game room by room name
	if (queries["room"][0] != "") {
		for _, game := range app.Games {
			if (game.Name == queries["room"][0]) {
				select_game = game	
			}
		}
		if (len(app.Games) > app.Configuration.Server.MaxRoom) {
			log.Println("[Error]: Server game room reach maximum number!")
			http.Error(w, "Server game room reach maximum number!", 400)
			return
		}
		// if the game room do not exist, then create new game room
		if (select_game == nil) {
			select_game = game.NewGame(queries["room"][0], 8192.0, 8192.0)
			// app.ControlLock.Lock()
			app.Games = append(app.Games, select_game)
			// app.ControlLock.Unlock()
		}
 	} else {
		// default select the first game instance
		select_game = app.Games[0]
	}
	// check repeatation of the player name
	// select_game.ControlLock.Lock()
	for _, ps := range select_game.Sessions {
		if (ps.Player.Attr.Name == queries["name"][0]) {
			log.Println("[Error]: Repeat player name!")
			http.Error(w, "Repeat player name!", 400)
			return
		}
	}
	// select_game.ControlLock.Unlock()
	// generate the player instance
	player := game.NewPlayer(queries["name"][0])
	// generate the player session
	session := game.NewSession(ws, player, select_game)
	// add the player session to the game
	select_game.JoinPlayer(session)

	// record some server information here
	(*app).Status.PlayerNumber += 1
	(*app).Status.RequestNumber += 1

	log.Printf("Player %s joined to game room %s", queries["name"][0], queries["room"][0])
	log.Printf("Server Status: %d room(s), %d player(s)", len((*app).Games), 1)
}

/**
 * <core>.staticHandler:
 * The function to handle the static file request
 *
 * @param {http.ResponseWriter} w 																							- the response writer of current request
 * @param {*http.Request} r																											- the current request
 *
 * @return {nil}
 */
func staticHandler (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/dist/" + r.URL.Path[1:])
}

/**
 * <App>.runServer:
 * The function in App to run server
 *
 * @return {nil}
 */
func (app App) runServer () {
	// handle the static file in url "/"
	http.HandleFunc("/", staticHandler)
	// handle the game websocket messaging in url "/game_ws"
	http.Handle("/game_ws", serverHandler { &app, gameWebsocketHandler })

	go func () {
		r := http.NewServeMux()
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		http.ListenAndServe(":3000", r)
	}()

	log.Println("run server on port", (*app.Configuration).Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", (*app.Configuration).Server.Host, (*app.Configuration).Server.Port), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}