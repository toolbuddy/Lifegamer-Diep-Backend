package game

import (
	"log"
)
// define the game frame rate
const frameRate = 3.0
// define game object struct
type Game struct {
	Name string
	Sessions []*PlayerSession
	MapInfo Map
	JoinChannel chan *PlayerSession
	Field *Size
}

func NewGame(name string) *Game {
	game := Game {
		Name: name,
		Sessions: []*PlayerSession {},
		JoinChannel: make(chan *PlayerSession),
		MapInfo: Map {
			Dieps: []*Diep {},
			Stuffs: []*Stuff {},
			Traps: []*Trap {},
		},
		Field: NewField(8192, 8192),
	}
	go game.runListen()
	go game.loop()
	return &game
}

// define the runListen function in Game struct pointer
func (g *Game) runListen () {
	for {
		// get the current join player session from channel
		p_sess := <- g.JoinChannel
		// append the player session to Sessions
		g.Sessions = append(g.Sessions, p_sess)
		log.Printf("Player %s has joined\n", p_sess.Player.Attr.Name)
	}
}
// define the loop function in Game struct pointer
func (g *Game) loop () {
	// TODO: collision detection
}
// define the JoinPlayer function in Game struct
func (g Game) JoinPlayer (session *PlayerSession) {
	// add the player session to channel
	g.JoinChannel <- session
}
