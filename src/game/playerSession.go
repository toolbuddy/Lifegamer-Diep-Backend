package game

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
	"log"
	"math"
	"sync"
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
	// "sort"
)

/**
 * PlayerView:
 * The struct to keep all instance on screen. It will be re-compute in every frame.
 *
 * @property {[]*GameObject} Dieps 			- all dieps in player views
 * @property {[]*GameObject} Stuffs			- all stuffs in player views
 * @property {[]*GameObject} Bullets		- all bullets in player views
 * @property {[]*GameObject} Traps			- all traps in player views
 */
type PlayerView struct {
	Dieps  []*GameObject // deips in view
	Stuffs []*GameObject // stuffs in view
	Bullets []*GameObject // bullet in view
	Traps []*GameObject // trap in view
}

/**
 * PlayerSession:
 * Every player owns one session, it will keep the neccessary reference to the player.
 *
 * @property {*websocket.Conn} Socket 	- websocket client connection instance
 * @property {*Game} Game 							- game room instance
 * @property {chan bool} MBus 					- the message channel between ping routine and all other routines
 * @property {bool} Alive 							- the status of the connection
 * @property {*Player} Player						- the player instance
 * @property {*PlayerView} View					- the view instance
 * @property {util.MoveDirection} Moving			- the current moving direction of player
 * @property {sync.Mutex} ControlLock		- the mutex lock to prevent from data race in routines
 */
type PlayerSession struct {
	Socket *websocket.Conn
	Game *Game
	MBus chan bool
	Alive bool
	Player *Player // player
	View PlayerView
	Moving util.MoveDirection
	ControlLock sync.Mutex
}

/**
 * CommandParams:
 * The definition of message param
 */
type CommandParams map[string]interface{}

/**
 * PlayerSessionCommand:
 * The definition of player message format
 *
 * @property {string} Method					 	- the action type of the player
 * @property {CommandParams} Params			- the option of the action
 */
type PlayerSessionCommand struct {
	Method string
	Params CommandParams
}

/**
 * <*PlayerSession>.receiver:
 * The function in PlayerSession to keep receiving message from client.
 */
func (ps *PlayerSession) receiver() {
	// keep read the player message
	for {
		if (!ps.Alive) {
			log.Println("receiver: not alive")
			return
		}
		_, command, err := ps.Socket.ReadMessage()
		if (err != nil) {
			break
		}
		var player_command PlayerSessionCommand = PlayerSessionCommand{}
		err = json.Unmarshal(command, &player_command)

		ps.serveCommand(player_command)
	}
}

/**
 * <*PlayerSession>.loop:
 * The function in PlayerSession to keep sending the player information in every frame.
 *
 * @return {nil}
 */
func (ps *PlayerSession) loop() {
	var stepDelay int32 = int32(1000 / ps.Game.Framerate)
	for {
		ps.ControlLock.Lock()
		if (!ps.Alive) {
			return
		}
		ps.ControlLock.Unlock()
		time.Sleep(time.Duration(stepDelay) * time.Millisecond)
		ps.sendPlayerState()
	}
}

/**
 * <*PlayerSession>.ping:
 * The function in PlayerSession to keep sending ping message to client verify connection.
 *
 * @return {nil}
 */
func (ps *PlayerSession) ping() {
	// use local variable here
	var alive bool = false
	for {
		if (!ps.Alive) {
			return
		}
		// send ping message to check connection alive first
		ps.sendPingMsg()
		// set alive false first
		alive = false
		// use chan bool and routine to count timeout
		timeout := make (chan bool, 1)
		go func () {
			time.Sleep(time.Second * 1)
			timeout <- true
		}()
		// catch channel message
		select {
			case  <- timeout:
				if (!alive) {
					log.Printf("Player %s disconnect", ps.Player.Attr.Name)
					ps.Game.Disconnect(ps.Player.Attr.Name)
					ps.Socket.Close()
					// lock the Alive attr in player session
					ps.ControlLock.Lock()
					ps.Alive = false
					// unlock the Alive attr in player session
					ps.ControlLock.Unlock()
					
					return
				}
			case <- ps.MBus:
				// set the alive vairable true if receive any message fron client
				alive = true
		}
		if (!alive) {
			return
		}
	}
}

/**
 * <*PlayerSession>.serveCommand:
 * The function in PlayerSession to deal with the message from client in PlayerSessionCommand format.
 *
 * @param {PlayerSessionCommand} command	- the command from client
 *
 * @return {nil}
 */
func (ps *PlayerSession) serveCommand(command PlayerSessionCommand) {
	// send the connection status through channel
	ps.MBus <- true
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
			ps.Shoot(command.Params["x"].(float64), command.Params["y"].(float64), int(command.Params["number"].(float64)))
			break
		case "evaluation":
			ps.Evaluation(command.Params["type"].(string))
			break
	}
}

/**
 * <*PlayerSession>.sendClientCommand:
 * The function in PlayerSession to send message to client in PlayerSessionCommand format.
 *
 * @param {PlayerSessionCommand} command	- the message sening to client
 *
 * @return {nil}
 */
func (ps *PlayerSession) sendClientCommand(command PlayerSessionCommand) {
	ps.ControlLock.Lock()
	message_b, _ := json.Marshal(command)
	err := ps.Socket.WriteMessage(websocket.TextMessage, message_b)
	ps.ControlLock.Unlock()
	if (err != nil) {
		ps.Socket.Close()
	}
}

/**
 * <*PlayerSession>.sendPlayerState:
 * The function in PlayerSession to send player status to client in PlayerSessionCommand format.
 *
 * @return {nil}
 */
func (ps *PlayerSession) sendPlayerState() {
	ps.ControlLock.Lock()
	// update the player view of all diep
	ps.updateView()
	ps.ControlLock.Unlock()
	// send all diep position to client
	ps.sendClientCommand(PlayerSessionCommand {
		Method: "playerSession",
		Params: CommandParams {
			"player": ps.Player,
			"dieps": ps.View.Dieps,
			"stuffs": ps.View.Stuffs,
			"traps": ps.View.Traps,
			"bullets": ps.View.Bullets,
		},
	})
	// log.Println("test")
}

/**
 * <*PlayerSession>.sendPingMsg:
 * The function in PlayerSession to send ping message to client in PlayerSessionCommand format.
 *
 * @return {nil}
 */
func (ps *PlayerSession) sendPingMsg() {
	ps.sendClientCommand(PlayerSessionCommand {
		Method: "ping",
		Params: CommandParams {},
	})
}

/**
 * <*PlayerSession>.updateView:
 * The function in PlayerSession to compute the view information in every frame.
 *
 * @return {nil}
 */
func (ps *PlayerSession) updateView () {
	// get the view width and height
	var vwL = math.Max(ps.Player.GameObject.Position.X - 1920 / 2, 0)
	var vwU = math.Min(ps.Player.GameObject.Position.X + 1920 / 2, ps.Game.Field.W)
	var vhL = math.Max(ps.Player.GameObject.Position.Y - 1080 / 2, 0)
	var vhU = math.Min(ps.Player.GameObject.Position.Y + 1080 / 2, ps.Game.Field.H)
	// empty all slice in player session
	ps.View.Dieps = []*GameObject {}
	ps.View.Stuffs = []*GameObject {}
	ps.View.Traps = []*GameObject {}
	ps.View.Bullets = []*GameObject {}
	// loop the map info and append the diep/stuff/trap in view
	for _, diep := range ps.Game.MapInfo.Dieps {
		if (diep.GameObject.Position.X >= vwL) && (diep.GameObject.Position.X <= vwU) &&
			(diep.GameObject.Position.Y >= vhL) && (diep.GameObject.Position.Y <= vhU) {
			ps.View.Dieps = append(ps.View.Dieps, &diep.GameObject)
		}
	}
	for _, stuff := range ps.Game.MapInfo.Stuffs {
		if (stuff.GameObject.Position.X >= vwL) && (stuff.GameObject.Position.X <= vwU) &&
			(stuff.GameObject.Position.Y >= vhL) && (stuff.GameObject.Position.Y <= vhU) {
			ps.View.Stuffs = append(ps.View.Stuffs, &stuff.GameObject)
		}
	}
	for _, trap := range ps.Game.MapInfo.Traps {
		if (trap.GameObject.Position.X >= vwL) && (trap.GameObject.Position.X <= vwU) &&
			(trap.GameObject.Position.Y >= vhL) && (trap.GameObject.Position.Y <= vhU) {
			ps.View.Traps = append(ps.View.Traps, &trap.GameObject)
		}
	}
	for _, bullet := range ps.Game.MapInfo.Bullets {
		if (bullet.GameObject.Position.X >= vwL) && (bullet.GameObject.Position.X <= vwU) &&
			(bullet.GameObject.Position.Y >= vhL) && (bullet.GameObject.Position.Y <= vhU) {
			ps.View.Bullets = append(ps.View.Bullets, &bullet.GameObject)
		}
	}
}

/**
 * <*PlayerSession>.Shoot:
 * The function in PlayerSession to shot.
 *
 * @param {float64} x										- the volume of angle project on x
 * @param {float64} y										- the volume of angle project on y
 * @param {int} number									- the bullet number in this shoot
 *
 * @return {nil}
 */
func (ps *PlayerSession) Shoot (x, y float64, number int) {
	// if there is cd time, then refuse shoot
	if (ps.Player.Attr.ShootCD >= 0) {
		ps.sendClientCommand(PlayerSessionCommand {
			Method: "shoot",
			Params: CommandParams {
				"message": "You can not shoot in CD time!",
			},
		})
		return
	}
	for i := 0; i < number; i++ {
		var new_bullet Bullet
		new_bullet.Position.X = ps.Player.Position.X
		new_bullet.Position.Y = ps.Player.Position.Y
		new_bullet.Velocity.X = x * float64(ps.Player.Status.BulletSpeed + 10) / ratio
		new_bullet.Velocity.Y = y * float64(ps.Player.Status.BulletSpeed + 10) / ratio
		new_bullet.Owner = ps.Player.Id
		new_bullet.Existence = (ps.Player.Status.BulletPenetration - 1) * 40 +  250
		ps.Game.MapInfo.Bullets = append(ps.Game.MapInfo.Bullets, &new_bullet)
	}
	// add shoot cd time
	ps.Player.Attr.ShootCD += 1000 / math.Log2(1.0 / float64(ps.Player.Status.BulletReload))
}

/**
 * <*PlayerSession>.Evaluation:
 * The function in PlayerSession to evaluation player diep.
 *
 * @param {string} type_str							- the chosen attribute of player in this evaluation
 *
 * @return {nil}
 */
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