package game

import (
	"github.com/gorilla/websocket"
	"log"
	"math"
	"sort"
	"time"
	"math/rand"
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
	"sync"
)

// define the game friction
const friction = 0.97
const ratio = 1.5

/**
 * Game:
 * The struct of game instance.
 *
 * @property {string} Name					 													- the unique identity between games
 * @property {[]*PlayerSessions} Sessions											- the slice of player sessions in the game
 * @property {Map} MapInfo																		- the map information
 * @property {chan *PlayerSessions} JoinChannel								- the channel of joining player
 * @property {*util.Size} Field																- the field information of the game
 * @property {float64} Framerate															- the framerate of the game
 * @property {*GameLogger} Logger															- the logger of the game
 */
 type Game struct {
	Name string
	Sessions []*PlayerSession
	MapInfo Map
	JoinChannel chan *PlayerSession
	Field *util.Size
	Framerate float64
	ControlLock sync.Mutex
	Logger *GameLogger
}

/**
 * <game>.NewGame:
 * The function to new a game instance.
 *
 * @param {string} name																				- the unique name of the game room
 * @param {float64} width																			- the width of the game field
 * @param {float64} height																		- the height of the game field
 *
 * @return {*Game}
 */
func NewGame(name string, width, height float64) *Game {
	game := Game {
		Name: name,
		Sessions: []*PlayerSession {},
		JoinChannel: make(chan *PlayerSession),
		MapInfo: Map {
			Dieps: []*Diep {},
			Stuffs: []*Stuff {},
			Traps: []*Trap {},
		},
		Field: &util.Size {
			W: width,
			H: height,
		},
		Framerate: 50.0,
		Logger: NewLogger(name),
	}
	go game.runListen()
	go game.loop()
	// generate the stuff randomly
	go game.generateStuff()
	return &game
}

/**
 * <game>.NewSession:
 * The function to new a player session.
 *
 * @param {*websocket.Conn} ws					- the websocket client instance
 * @param {*Player} player							- the player instance
 * @param {*Game} game									- the game instance
 *
 * @return {*PlayerSession}
 */
 func NewSession(ws *websocket.Conn, player *Player, game *Game) *PlayerSession {
	// init the session
	ps := PlayerSession {
		Socket: ws,
		Player: player,
		Game: game,
		MBus: make (chan bool, 1),
		Alive: true,
	}
	// log the connection
	game.Logger.establishConnection(ws.RemoteAddr().String(), player.Attr.Name, game.Name, int(len(game.Sessions)) + 1)
	// parallel execute receiver, loop and ping function
	go ps.receiver()
	go ps.loop()
	go ps.ping()
	return &ps
}

/**
 * <game>.NewPlayer:
 * The function to new a player.
 *
 * @param {stirng} name									- the name of the player
 *
 * @return {*Player}
 */
func NewPlayer (name string) *Player {
	uuid, _ := util.NewUUID()
	var new_player = Player {
		GameObject: GameObject {
			Id: uuid,
			Position: util.Point {
				X: rand.Float64() * 1023,
				Y: rand.Float64() * 1023,
			},
			Mass: 1.0,
			Radius: 50.0,
			Velocity: util.VelocityFormat {
				X: 0.0,
				Y: 0.0,
			},
			Acceleration: util.AccelerationFormat {
				Up: 0.0,
				Down: 0.0,
				Left: 0.0,
				Right: 0.0,
			},
		},
		Attr: PlayerAttribute {
			Name: name,
			CreatedAt: time.Now(),
			Score: 0,
			Level: 1,
			EXP: 0,
			HP: 100,
			ShootCD: 0,
		},
		Status: PlayerStatus {
			MaxHP: 100,
			HPRegeneration: 1,
			MoveSpeed: 1,
			BulletSpeed: 1,
			BulletPenetration: 1,
			BulletReload: 1,
			BulletDamage: 1,
			BodyDamage: 1,
		},
	}
	return &new_player
}

/**
 * <game>.GetObjectId:
 * The function to get game object id through GameObjectInterface.
 *
 * @return {string}
 */
 func GetObjectId (m GameObjectInterface) string {
	return m.GetId()
}

/**
 * <*Game>.runListen:
 * The function in Game to listen the join request from server.
 *
 * @return {nil}
 */
func (g *Game) runListen () {
	for {
		// get the current join player session from channel
		p_sess := <- g.JoinChannel
		g.ControlLock.Lock()
		// append the player session to Sessions
		g.Sessions = append(g.Sessions, p_sess)
		g.ControlLock.Unlock()
		log.Printf("Player %s has joined\n", p_sess.Player.Attr.Name)
	}
}

/**
 * <*Game>.JoinPlayer:
 * The function in Game to send session to the channel.
 *
 * @return {string}
 */
 func (g *Game) JoinPlayer (session *PlayerSession) {
	// add the player session to channel
	g.JoinChannel <- session
}

/**
 * <*Game>.Disconnect:
 * The function in Game to remove player session with disconnection.
 *
 * @property {string} player_name				- the target player name
 *
 * @return {nil}
 */
 func (g *Game) Disconnect (player_name string) {
	 // log the disconnection
	g.Logger.closeConnection(player_name, g.Name, int(len(g.Sessions)) - 1)
	// remove the player session from the game
	var index int = sort.Search(len(g.Sessions), func (i int) bool {
		return g.Sessions[i].Player.Attr.Name == player_name
	})
	g.ControlLock.Lock()
	g.Sessions = append(g.Sessions[:index], g.Sessions[index+1:]...)
	g.ControlLock.Unlock()
}

/**
 * <*Game>.loop:
 * The function in Game to keep computing the all movement of the item in the game.
 *
 * @return {nil}
 */
func (g *Game) loop () {
	for {
		time.Sleep(time.Duration(1000.0 / g.Framerate) * time.Millisecond)
		g.ControlLock.Lock()
		// update the player movement
		g.updatePhysicItems()
		g.ControlLock.Unlock()
		// detect player & player collision
		g.detectDeipCollision()
		// detect player & stuff collision
		g.detectStuffCollision()
		// detect player & trap collision
		g.detectTrapCollision()
		// detect player & bullet collision
		g.detectBulletCollision()
		// deal all collision
		g.dealWithCollisions()
	}
}

/**
 * <*Game>.updatePhysicItems:
 * The function in Game to computing all game object movement.
 *
 * @return {nil}
 */
func (g *Game) updatePhysicItems() {
	// update the player movement
	for _, ps := range g.Sessions {
		ps.ControlLock.Lock()
		// update the player acceleration
		var new_acceleration util.AccelerationFormat
		if (ps.Moving.Up) {
			new_acceleration.Up = math.Min(ps.Player.GameObject.Acceleration.Up + (float64(ps.Player.Status.MoveSpeed)) * friction / g.Framerate, float64(ps.Player.Status.MoveSpeed))
		} else {
			new_acceleration.Up = math.Max(ps.Player.GameObject.Acceleration.Up * friction, 0) / g.Framerate
		}
		if (ps.Moving.Down) {
			new_acceleration.Down = math.Min(ps.Player.GameObject.Acceleration.Down + (float64(ps.Player.Status.MoveSpeed)) * friction / g.Framerate, float64(ps.Player.Status.MoveSpeed))
		} else {
			new_acceleration.Down = math.Max(ps.Player.GameObject.Acceleration.Down * friction, 0) / g.Framerate
		}
		if (ps.Moving.Left) {
			new_acceleration.Left = math.Min(ps.Player.GameObject.Acceleration.Left + (float64(ps.Player.Status.MoveSpeed)) * friction / g.Framerate, float64(ps.Player.Status.MoveSpeed))
		} else {
			new_acceleration.Left = math.Max(ps.Player.GameObject.Acceleration.Left * friction, 0) / g.Framerate
		}
		if (ps.Moving.Right) {
			new_acceleration.Right = math.Min(ps.Player.GameObject.Acceleration.Right + (float64(ps.Player.Status.MoveSpeed)) * friction / g.Framerate, float64(ps.Player.Status.MoveSpeed))
		} else {
			new_acceleration.Right = math.Max(ps.Player.GameObject.Acceleration.Right * friction, 0) / g.Framerate
		}
		ps.Player.GameObject.Acceleration = new_acceleration
		// update the player velocity
		ps.Player.GameObject.Velocity.X = math.Max(math.Min(ps.Player.GameObject.Velocity.X - ps.Player.GameObject.Acceleration.Left +
			ps.Player.GameObject.Acceleration.Right, float64(ps.Player.Status.MoveSpeed + 10) / ratio), float64(ps.Player.Status.MoveSpeed + 10) * (-1.0) / ratio) * friction
		ps.Player.GameObject.Velocity.Y = math.Max(math.Min(ps.Player.GameObject.Velocity.Y - ps.Player.GameObject.Acceleration.Up +
				ps.Player.GameObject.Acceleration.Down, float64(ps.Player.Status.MoveSpeed + 10) / ratio), float64(ps.Player.Status.MoveSpeed + 10) * (-1.0) / ratio) * friction

		// update the player location
		ps.Player.GameObject.Position.X = math.Max(math.Min(ps.Player.GameObject.Position.X + ps.Player.GameObject.Velocity.X / g.Framerate, g.Field.W), 0)
		ps.Player.GameObject.Position.Y = math.Max(math.Min(ps.Player.GameObject.Position.Y + ps.Player.GameObject.Velocity.Y / g.Framerate, g.Field.H), 0)
		
		// update the player shoot CD time
		ps.Player.Attr.ShootCD = math.Max(ps.Player.Attr.ShootCD - 1, 0)

		ps.ControlLock.Unlock()
	}
	// update the bullet movement
	for i, bullet := range g.MapInfo.Bullets {
		bullet.GameObject.Position.X = math.Max(math.Min(bullet.GameObject.Position.X + bullet.GameObject.Velocity.X / g.Framerate, g.Field.W), 0)
		bullet.GameObject.Position.Y = math.Max(math.Min(bullet.GameObject.Position.Y + bullet.GameObject.Velocity.Y / g.Framerate, g.Field.H), 0)
		// if the bullet collide with wall
		if (bullet.GameObject.Position.X >= g.Field.W) || (bullet.GameObject.Position.X <= 0) ||
			(bullet.GameObject.Position.Y >= g.Field.H) || (bullet.GameObject.Position.Y <= 0) {
			// remove the bullet
			g.MapInfo.Bullets = append(g.MapInfo.Bullets[:i], g.MapInfo.Bullets[i+1:]...)
		}
		// count for the bullet existence
		bullet.Existence--;
		if (bullet.Existence <= 0) {
			// remove the bullet
			g.MapInfo.Bullets = append(g.MapInfo.Bullets[:i], g.MapInfo.Bullets[i+1:]...)
		}
	}
	// update the stuff movement
	for _, stuff := range g.MapInfo.Stuffs {
		// update the stuff acceleration
		var new_acceleration util.AccelerationFormat
		new_acceleration.Up = math.Max(stuff.GameObject.Acceleration.Up * friction, 0) / g.Framerate
		new_acceleration.Down = math.Max(stuff.GameObject.Acceleration.Down * friction, 0) / g.Framerate
		new_acceleration.Left = math.Max(stuff.GameObject.Acceleration.Left * friction, 0) / g.Framerate
		new_acceleration.Right = math.Max(stuff.GameObject.Acceleration.Right * friction, 0) / g.Framerate
		stuff.GameObject.Acceleration = new_acceleration
		// update the player velocity
		stuff.GameObject.Velocity.X = (stuff.GameObject.Velocity.X - stuff.GameObject.Acceleration.Left +
			stuff.GameObject.Acceleration.Right) * friction
		stuff.GameObject.Velocity.Y = (stuff.GameObject.Velocity.Y - stuff.GameObject.Acceleration.Up +
				stuff.GameObject.Acceleration.Down) * friction

		// update the player location
		stuff.GameObject.Position.X = math.Max(math.Min(stuff.GameObject.Position.X + stuff.GameObject.Velocity.X / g.Framerate, g.Field.W), 0)
		stuff.GameObject.Position.Y = math.Max(math.Min(stuff.GameObject.Position.Y + stuff.GameObject.Velocity.Y / g.Framerate, g.Field.H), 0)
	}
}

/**
 * <*Game>.detectDeipCollision:
 * The function in Game to detect if there is collision between two diep.
 *
 * @return {nil}
 */
func (g *Game) detectDeipCollision() {
	for _, diep_a := range g.MapInfo.Dieps {
		for _, diep_b := range g.MapInfo.Dieps {
			if (diep_a.GameObject.Id == diep_b.GameObject.Id) {
				continue
			}
			var diff_x = math.Abs(diep_a.GameObject.Position.X - diep_b.GameObject.Position.X)
			var diff_y = math.Abs(diep_a.GameObject.Position.Y - diep_b.GameObject.Position.Y)
			var check = false
			if (math.Sqrt(math.Pow(diff_x, 2) + math.Pow(diff_y, 2)) <= diep_a.Radius + diep_b.Radius) {
				check = true
			}
			
			// collision happend, then add the acceleration in opposite direction
			if (check) {
				pair_exist := sort.Search(len(g.MapInfo.Collisions), func (i int) bool {
					targetId := GetObjectId(g.MapInfo.Collisions[i].object_b)
					return (g.MapInfo.Collisions[i].object_a.GameObject.Id == diep_a.GameObject.Id && targetId == diep_b.GameObject.Id) || 
						(g.MapInfo.Collisions[i].object_a.GameObject.Id == diep_b.GameObject.Id && targetId == diep_a.GameObject.Id)
				})
				if (pair_exist >= len(g.MapInfo.Collisions)) {
					continue
				}
				// prevent the collision detect twice by add the pair to the list unrepeatly
				g.MapInfo.Collisions = append(g.MapInfo.Collisions, CollisionDetection {
					object_a: diep_a,
					object_b: diep_b,
				})
			}
		}
	}
}

/**
 * <*Game>.detectStuffCollision:
 * The function in Game to detect if there is collision between diep and stuff.
 *
 * @return {nil}
 */
func (g *Game) detectStuffCollision () {
	for _, diep := range g.MapInfo.Dieps {
		for _, stuff := range g.MapInfo.Stuffs {
			var diff_x = math.Abs(diep.GameObject.Position.X - stuff.GameObject.Position.X)
			var diff_y = math.Abs(diep.GameObject.Position.Y - stuff.GameObject.Position.Y)
			var check = false
			if (math.Sqrt(math.Pow(diff_x, 2) + math.Pow(diff_y, 2)) <= diep.Radius + stuff.Radius) {
				check = true
			}
			
			// collision happend, then add the acceleration in opposite direction
			if (check) {
				pair_exist := sort.Search(len(g.MapInfo.Collisions), func (i int) bool {
					targetId := GetObjectId(g.MapInfo.Collisions[i].object_b)
					return (g.MapInfo.Collisions[i].object_a.GameObject.Id == diep.Id && targetId == stuff.Id) || 
						(g.MapInfo.Collisions[i].object_a.GameObject.Id == stuff.Id && targetId == diep.Id)
				})
				if (pair_exist >= len(g.MapInfo.Collisions)) {
					continue
				}
				// prevent the collision detect twice by add the pair to the list unrepeatly
				g.MapInfo.Collisions = append(g.MapInfo.Collisions, CollisionDetection {
					object_a: diep,
					object_b: stuff,
				})
			}
		}
	}
}

/**
 * <*Game>.detectTrapCollision:
 * The function in Game to detect if there is collision between diep and trap.
 *
 * @return {nil}
 */
func (g *Game) detectTrapCollision () {
	for _, diep := range g.MapInfo.Dieps {
		for _, trap := range g.MapInfo.Traps {
			var diff_x = math.Abs(diep.GameObject.Position.X - trap.GameObject.Position.Y)
			var diff_y = math.Abs(diep.GameObject.Position.X - trap.GameObject.Position.Y)
			var check = false
			if (math.Sqrt(math.Pow(diff_x, 2) + math.Pow(diff_y, 2)) <= diep.Radius + trap.Radius) {
				check = true
			}
			
			// collision happend, then add the acceleration in opposite direction
			if (check) {

				pair_exist := sort.Search(len(g.MapInfo.Collisions), func (i int) bool {
					targetId := GetObjectId(g.MapInfo.Collisions[i].object_b)
					return (g.MapInfo.Collisions[i].object_a.GameObject.Id == diep.Id && targetId == trap.Id) || 
						(g.MapInfo.Collisions[i].object_a.GameObject.Id == trap.Id && targetId == diep.Id)
				})
				if (pair_exist >= len(g.MapInfo.Collisions)) {
					continue
				}
				// prevent the collision detect twice by add the pair to the list unrepeatly
				g.MapInfo.Collisions = append(g.MapInfo.Collisions, CollisionDetection {
					object_a: diep,
					object_b: trap,
				})
			}
		}
	}
}

/**
 * <*Game>.detectBulletCollision:
 * The function in Game to detect if there is collision between diep and bullet.
 *
 * @return {nil}
 */
func (g *Game) detectBulletCollision () {
	for _, diep := range g.MapInfo.Dieps {
		for _, bullet := range g.MapInfo.Bullets {
			var diff_x = math.Abs(diep.GameObject.Position.X - bullet.GameObject.Position.Y)
			var diff_y = math.Abs(diep.GameObject.Position.X - bullet.GameObject.Position.Y)
			var check = false
			if (math.Sqrt(math.Pow(diff_x, 2) + math.Pow(diff_y, 2)) <= diep.Radius + bullet.Radius) {
				check = true
			}
			
			// collision happend, then add the acceleration in opposite direction
			if (check) {

				pair_exist := sort.Search(len(g.MapInfo.Collisions), func (i int) bool {
					targetId := GetObjectId(g.MapInfo.Collisions[i].object_b)
					return (g.MapInfo.Collisions[i].object_a.GameObject.Id == diep.Id && targetId == bullet.Id) || 
						(g.MapInfo.Collisions[i].object_a.GameObject.Id == bullet.Id && targetId == diep.Id)
				})
				if (pair_exist >= len(g.MapInfo.Collisions)) {
					continue
				}
				// prevent the collision detect twice by add the pair to the list unrepeatly
				g.MapInfo.Collisions = append(g.MapInfo.Collisions, CollisionDetection {
					object_a: diep,
					object_b: bullet,
				})
			}
		}
	}
}

/**
 * <*Game>.dealWithCollisions:
 * The function in Game to apply all collision effect.
 *
 * @return {nil}
 */
func (g *Game) dealWithCollisions () {
	for _, collision := range g.MapInfo.Collisions {
		// object_a must bee diep, then just get the player_a session
		player_session_a_index := sort.Search(len(g.Sessions), func (i int) bool {
			return g.Sessions[i].Player.GameObject.Id == collision.object_a.GameObject.Id
		})
		var player_session_a = g.Sessions[player_session_a_index]

		switch collision.object_b.(type) {
			case *Diep:
				target := collision.object_b.(*Diep)
				player_session_b_index := sort.Search(len(g.Sessions), func (i int) bool {
					return g.Sessions[i].Player.GameObject.Id == target.GameObject.Id
				})
				var player_session_b = g.Sessions[player_session_b_index]
				// update the acceleration of two player
				var new_acceleration_a = util.AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up + player_session_b.Player.GameObject.Acceleration.Up
				new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down + player_session_b.Player.GameObject.Acceleration.Down
				new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left + player_session_b.Player.GameObject.Acceleration.Left
				new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right + player_session_b.Player.GameObject.Acceleration.Right

				var new_acceleration_b = util.AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_b.Up = player_session_b.Player.GameObject.Acceleration.Up + player_session_a.Player.GameObject.Acceleration.Up
				new_acceleration_b.Down = player_session_b.Player.GameObject.Acceleration.Down + player_session_a.Player.GameObject.Acceleration.Down
				new_acceleration_b.Left = player_session_b.Player.GameObject.Acceleration.Left + player_session_a.Player.GameObject.Acceleration.Left
				new_acceleration_b.Right = player_session_b.Player.GameObject.Acceleration.Right + player_session_a.Player.GameObject.Acceleration.Right
				
				player_session_a.Player.GameObject.Acceleration = new_acceleration_a
				player_session_b.Player.GameObject.Acceleration = new_acceleration_b
				
				// give the collision damage
				player_session_a.Player.Attr.HP -= float64(player_session_b.Player.Status.BodyDamage) * 5.0
				player_session_b.Player.Attr.HP -= float64(player_session_a.Player.Status.BodyDamage) * 5.0
				// deal with the dead
				if (player_session_a.Player.Attr.HP <= 0) {
					// log the dead message
					g.Logger.deadMessage(player_session_a.Player.GameObject.Id, player_session_b.Player.GameObject.Id)
					// send the dead message first
					player_session_a.sendClientCommand(PlayerSessionCommand {
						Method: "playerDead",
						Params: CommandParams {},
					})
					player_session_a.ControlLock.Lock()
					player_session_a.Alive = false
					player_session_a.ControlLock.Unlock()
				}
				if (player_session_b.Player.Attr.HP <= 0) {
					// log the dead message
					g.Logger.deadMessage(player_session_b.Player.GameObject.Id, player_session_a.Player.GameObject.Id)
					// send the dead message first
					player_session_b.sendClientCommand(PlayerSessionCommand {
						Method: "playerDead",
						Params: CommandParams {},
					})
					player_session_b.ControlLock.Lock()
					player_session_b.Alive = false
					player_session_b.ControlLock.Unlock()
				}
				break;
			case *Stuff:
				target := collision.object_b.(*Stuff)
				stuff_index := sort.Search(len(g.MapInfo.Stuffs), func (i int) bool {
					return g.MapInfo.Stuffs[i].GameObject.Id == target.GameObject.Id
				})
				var stuff = g.MapInfo.Stuffs[stuff_index]
				// update the acceleration of two player
				var new_acceleration_a = util.AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up + stuff.Acceleration.Up
				new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down + stuff.Acceleration.Down
				new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left + stuff.Acceleration.Left
				new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right + stuff.Acceleration.Right

				var new_acceleration_s = util.AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_s.Up = stuff.Acceleration.Up + player_session_a.Player.GameObject.Acceleration.Up
				new_acceleration_s.Down = stuff.Acceleration.Down + player_session_a.Player.GameObject.Acceleration.Down
				new_acceleration_s.Left = stuff.Acceleration.Left + player_session_a.Player.GameObject.Acceleration.Left
				new_acceleration_s.Right = stuff.Acceleration.Right + player_session_a.Player.GameObject.Acceleration.Right
				
				player_session_a.Player.GameObject.Acceleration = new_acceleration_a
				stuff.Acceleration = new_acceleration_s
				
				// give the collision damage
				player_session_a.Player.Attr.HP -= float64(stuff.Attr.BodyDamage) * 5.0
				stuff.Attr.HP -= float64(player_session_a.Player.Status.BodyDamage) * 5.0
				// deal with the dead
				if (player_session_a.Player.Attr.HP <= 0) {
					// log the dead message
					g.Logger.deadMessage(player_session_a.Player.GameObject.Id, stuff.GameObject.Id)
					// send the dead message first
					player_session_a.sendClientCommand(PlayerSessionCommand {
						Method: "playerDead",
						Params: CommandParams {},
					})
					player_session_a.ControlLock.Lock()
					player_session_a.Alive = false
					player_session_a.ControlLock.Unlock()
				}
				if (stuff.Attr.HP <= 0) {
					// log the dead message
					g.Logger.deadMessage(stuff.GameObject.Id, player_session_a.Player.GameObject.Id)
					player_session_a.Player.GainEXP(stuff.Attr.EXP)
					// remove the stuff
					g.MapInfo.Stuffs = append(g.MapInfo.Stuffs[:stuff_index], g.MapInfo.Stuffs[stuff_index+1:]...)
				}
				break;
			case *Trap:
				target := collision.object_b.(*Trap)
				trap_index := sort.Search(len(g.MapInfo.Traps), func (i int) bool {
					return g.MapInfo.Traps[i].GameObject.Id == target.GameObject.Id
				})
				var trap = g.MapInfo.Traps[trap_index]
				// update the acceleration of two player
				var new_acceleration_a = util.AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				if (player_session_a.Player.GameObject.Acceleration.Left > 0.0) || (player_session_a.Player.GameObject.Acceleration.Right > 0.0) {
					new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up
					new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down
					new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left * (-1.0)
					new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right * (-1.0)
				} else {
					new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up * (-1.0)
					new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down * (-1.0)
					new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left
					new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right
				}

				player_session_a.Player.GameObject.Acceleration = new_acceleration_a
				
				// give the collision damage
				player_session_a.Player.Attr.HP -= float64(trap.Attr.BodyDamage) * 5.0
				// deal with the dead
				if (player_session_a.Player.Attr.HP <= 0) {
					// log the dead message
					g.Logger.deadMessage(player_session_a.Player.GameObject.Id, trap.GameObject.Id)
					// send the dead message first
					player_session_a.sendClientCommand(PlayerSessionCommand {
						Method: "playerDead",
						Params: CommandParams {},
					})
					player_session_a.ControlLock.Lock()
					player_session_a.Alive = false
					player_session_a.ControlLock.Unlock()
				}
				if (trap.Attr.HP <= 0) {
					// log the dead message
					g.Logger.deadMessage(trap.GameObject.Id, player_session_a.Player.GameObject.Id)
					// remove the trap
					g.MapInfo.Traps = append(g.MapInfo.Traps[:trap_index], g.MapInfo.Traps[trap_index+1:]...)
				}
				break;
			case *Bullet:
				target := collision.object_b.(*Bullet)
				bullet_index := sort.Search(len(g.MapInfo.Bullets), func (i int) bool {
					return g.MapInfo.Bullets[i].GameObject.Id == target.GameObject.Id
				})
				var bullet = g.MapInfo.Bullets[bullet_index]
				
				// give the collision damage
				if (bullet.Owner != player_session_a.Player.Attr.Name) {
					player_session_a.Player.Attr.HP -= float64(bullet.Damage) * 5.0
					if (player_session_a.Player.Attr.HP <= 0) {
						// log the dead message
						g.Logger.deadMessage(player_session_a.Player.GameObject.Id, bullet.GameObject.Id)
						// log the dead message
						g.Logger.deadMessage(player_session_a.Player.GameObject.Id, bullet.GameObject.Id)
						// send the dead message first
						player_session_a.sendClientCommand(PlayerSessionCommand {
							Method: "playerDead",
							Params: CommandParams {},
						})
						player_session_a.ControlLock.Lock()
						player_session_a.Alive = false
						player_session_a.ControlLock.Unlock()
					}
					// remove the bullet
					g.MapInfo.Bullets = append(g.MapInfo.Bullets[:bullet_index], g.MapInfo.Bullets[bullet_index+1:]...)
				}
				break;
		}
	}
}
