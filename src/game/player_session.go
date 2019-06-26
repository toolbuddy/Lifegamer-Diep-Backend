package game

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
	"log"
)

type MoveDirection struct {
	Up bool
	Down bool
	Left bool
	Right bool
}

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
	Moving MoveDirection
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
	ps := PlayerSession {
		Socket: ws,
		Player: player,
		Game: game,
	}
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
			log.Fatal(err)
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
	log.Println(command.Params["Value"])
	// define player method with correspond action
	switch command.Method {
		case "moveUp":
			ps.Moving.Up = command.Params["value"].(bool)
			break
		case "moveDown":
			ps.Moving.Down = command.Params["value"].(bool)
			break
		case "moveLeft":
			ps.Moving.Left = command.Params["value"].(bool)
			break
		case "moveRight":
			ps.Moving.Right = command.Params["value"].(bool)
			break
		case "shoot":
			break
		case "evaluation":
			ps.Evaluation(command.Params["type"].(string))
			break
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

// define the updateViewInfo function in PlayerSession Pointer
func (ps *PlayerSession) updateViewInfo (m *Map) {
	// get the view width and height
	var vwL = ps.Player.GameObject.Position.X - 1920 / 2
	if (vwL < 0) {
		vwL = 0
	}
	var vwU = ps.Player.GameObject.Position.X + 1920 / 2
	if (vwU > 8192) {
		vwU = 8192
	}
	var vhL = ps.Player.GameObject.Position.Y - 1080 / 2
	if (vhL < 0) {
		vhL = 0
	}
	var vhU = ps.Player.GameObject.Position.Y + 1080 / 2
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
		if (diep.GameObject.Position.X >= vwL) && (diep.GameObject.Position.X <= vwU) &&
			(diep.GameObject.Position.Y >= vhL) && (diep.GameObject.Position.Y <= vhU) {
			ps.Dieps = append(ps.Dieps, diep)
		}
	}
	for _, stuff := range m.Stuffs {
		if (stuff.GameObject.Position.X >= vwL) && (stuff.GameObject.Position.X <= vwU) &&
			(stuff.GameObject.Position.Y >= vhL) && (stuff.GameObject.Position.Y <= vhU) {
			ps.Stuffs = append(ps.Stuffs, stuff)
		}
	}
	for _, trap := range m.Traps {
		if (trap.GameObject.Position.X >= vwL) && (trap.GameObject.Position.X <= vwU) &&
			(trap.GameObject.Position.Y >= vhL) && (trap.GameObject.Position.Y <= vhU) {
			ps.Traps = append(ps.Traps, trap)
		}
	}
	for _, bullet := range m.Bullets {
		if (bullet.GameObject.Position.X >= vwL) && (bullet.GameObject.Position.X <= vwU) &&
			(bullet.GameObject.Position.Y >= vhL) && (bullet.GameObject.Position.Y <= vhU) {
			ps.Bullets = append(ps.Bullets, bullet)
		}
	}
}

func (ps *PlayerSession) Evaluation (type_str string) {
	switch type_str {
		case "MaxHP":
			ps.Player.Status.MaxHP++
			break
		case "HPRegeneration":
			ps.Player.Status.HPRegeneration++
			break
		case "MoveSpeed":
			ps.Player.Status.MoveSpeed++
			break
		case "BulletSpeed":
			ps.Player.Status.BulletSpeed++
			break
		case "BulletPenetration":
			ps.Player.Status.BulletPenetration++
			break
		case "BulletReload":
			ps.Player.Status.BulletReload++
			break
		case "BulletDamage":
			ps.Player.Status.BulletDamage++
			break
		case "BodyDamage":
			ps.Player.Status.BodyDamage++
			break
	}
}