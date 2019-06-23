package game

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"time"
)

// define PlayerSession struct
type PlayerSession struct {
	Socket *websocket.Conn
	Game *Game
	Player *Player // player
	Dieps  []*Diep // deips in view
	Stuffs []*Stuff // stuffs in view
	Bullets []*Bullet // bullet in view
	Traps []*Trap // trap in view
	Direction Point // shoot angle
}
// define CommandpParams struct
type CommandParams map[string]interface{}
// define PlayerSessionCommand struct
type PlayerSessionCommand struct {
	Method string
	Params CommandParams
}

// define the NewSession function in game package
func NewSession(ws *websocket.Conn, player *Player, game *Game) *PlayerSession {
	// init the session
	ps := PlayerSession{Socket: ws, Player: player, Game: game}
	// parallel execute receiver and loop function
	go ps.receiver()
	go ps.loop(&game.MapInfo)
	return &ps
}
// define the receiver function in PlayerSession pointer
func (ps *PlayerSession) receiver() {
	// keep read the player message
	for {
		_, command, err := ps.Socket.ReadMessage()
		if err != nil {
			break
		}
		var player_command PlayerSessionCommand = PlayerSessionCommand{}
		err = json.Unmarshal(command, &player_command)

		ps.serverCommand(player_command)
	}
	ps.Socket.Close()
}
// define the serverCommand function in PlayerSession pointer
func (ps *PlayerSession) serverCommand(command PlayerSessionCommand) {
	// define player method with correspond action
	switch command.Method {
		case "updateDirection":
		case "moveUp":
		case "moveDown":
		case "moveLeft":
		case "moveRight":
		case "stop":
		case "shoot":
		case "evaluation":
	}
}

func (ps *PlayerSession) sendClientCommand(command PlayerSessionCommand) {
	message_b, _ := json.Marshal(command)
	err := ps.Socket.WriteMessage(websocket.TextMessage, message_b)
	if err != nil {
		ps.Socket.Close()
	}
}

func (ps *PlayerSession) loop(m *Map) {
	var stepDelay int32 = 20
	for {
		time.Sleep(time.Duration(stepDelay) * time.Millisecond)
		ps.sendPlayerState(m)
	}
}

func (ps *PlayerSession) sendPlayerState(m *Map) {
	// update the player view of all diep
	ps.updateViewInfo(m)
	// send all diep position to client
	ps.sendClientCommand(PlayerSessionCommand {
		Method: "playerSession",
		Params: CommandParams {
			"player": ps.Player,
			"dieps": ps.Dieps,
			"stuffs": ps.Stuffs,
			"traps": ps.Traps,
		},
	})
}

func (ps *PlayerSession) updateDirection(command PlayerSessionCommand) {
	xDif := math.Cos(command.Params["directionAngelRad"].(float64)) * ps.getSpeed()
	yDif := math.Sin(command.Params["directionAngelRad"].(float64)) * ps.getSpeed()
	log.Printf("update Direction : x : %f, y: %f  \n", xDif, yDif)
	ps.Direction.X = float64(xDif)
	ps.Direction.Y = float64(yDif)
}

// define the getSpeed function in PlayerSession Pointer
func (ps *PlayerSession) getSpeed() float64 {
	return float64(ps.Player.Status.MoveSpeed) * 0.5 + 1
}

// define the updateViewInfo function in PlayerSession Pointer
func (ps *PlayerSession) updateViewInfo (m *Map) {
	// get the view width and height
	var vwL = ps.Player.Attr.Position.X - 1920 / 2
	if (vwL < 0) {
		vwL = 0
	}
	var vwU = ps.Player.Attr.Position.X + 1920 / 2
	if (vwU > 8192) {
		vwU = 8192
	}
	var vhL = ps.Player.Attr.Position.Y - 1080 / 2
	if (vhL < 0) {
		vhL = 0
	}
	var vhU = ps.Player.Attr.Position.Y + 1080 / 2
	if (vhU > 8192) {
		vhU = 8192
	}
	// empty all slice in player session
	ps.Dieps = []*Diep {}
	ps.Stuffs = []*Stuff {}
	ps.Traps = []*Trap {}
	ps.Bullets = []*Bullet {}
	// loop the map info and append the diep/stuff/trap in view
	for _, diep := range m.Dieps {
		if (diep.Position.X >= vwL) && (diep.Position.X <= vwU) &&
			(diep.Position.Y >= vhL) && (diep.Position.Y <= vhU) {
			ps.Dieps = append(ps.Dieps, diep)
		}
	}
	for _, stuff := range m.Stuffs {
		if (stuff.Position.X >= vwL) && (stuff.Position.X <= vwU) &&
			(stuff.Position.Y >= vhL) && (stuff.Position.Y <= vhU) {
			ps.Stuffs = append(ps.Stuffs, stuff)
		}
	}
	for _, trap := range m.Traps {
		if (trap.Position.X >= vwL) && (trap.Position.X <= vwU) &&
			(trap.Position.Y >= vhL) && (trap.Position.Y <= vhU) {
			ps.Traps = append(ps.Traps, trap)
		}
	}
	for _, bullet := range m.Bullets {
		if (bullet.Position.X >= vwL) && (bullet.Position.X <= vwU) &&
			(bullet.Position.Y >= vhL) && (bullet.Position.Y <= vhU) {
			ps.Bullets = append(ps.Bullets, bullet)
		}
	}
}