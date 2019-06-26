package game

import (
	"log"
	"math"
	"sort"
	"time"
)

// define the game friction
const friction = 0.97
const ratio = 1.5
// define game object struct
type Game struct {
	Name string
	Sessions []*PlayerSession
	MapInfo Map
	JoinChannel chan *PlayerSession
	Field *Size
	Framerate float64
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
		Framerate: 50.0,
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
	for {
		time.Sleep(time.Duration(1000.0 / g.Framerate) * time.Millisecond)
		// update the player movement
		g.updatePhysicItems()
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

func (g *Game) updatePhysicItems() {
	for _, ps := range g.Sessions {
		// update the player acceleration
		var new_acceleration AccelerationFormat
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
		log.Println("acceleration", ps.Player.GameObject.Position)
	}
}

// define the detectDiepCollision in game package
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
				var new_acceleration_a = AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up + player_session_b.Player.GameObject.Acceleration.Up
				new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down + player_session_b.Player.GameObject.Acceleration.Down
				new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left + player_session_b.Player.GameObject.Acceleration.Left
				new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right + player_session_b.Player.GameObject.Acceleration.Right

				var new_acceleration_b = AccelerationFormat {
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
				break;
			case *Stuff:
				target := collision.object_b.(*Stuff)
				stuff_index := sort.Search(len(g.MapInfo.Stuffs), func (i int) bool {
					return g.MapInfo.Stuffs[i].GameObject.Id == target.GameObject.Id
				})
				var stuff = g.MapInfo.Stuffs[stuff_index]
				// update the acceleration of two player
				var new_acceleration_a = AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up + stuff.Acceleration.Up
				new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down + stuff.Acceleration.Down
				new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left + stuff.Acceleration.Left
				new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right + stuff.Acceleration.Right

				var new_acceleration_s = AccelerationFormat {
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
				break;
			case *Trap:
				target := collision.object_b.(*Trap)
				trap_index := sort.Search(len(g.MapInfo.Traps), func (i int) bool {
					return g.MapInfo.Traps[i].GameObject.Id == target.GameObject.Id
				})
				var trap = g.MapInfo.Traps[trap_index]
				// update the acceleration of two player
				var new_acceleration_a = AccelerationFormat {
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
				break;
			case *Bullet:
				target := collision.object_b.(*Bullet)
				bullet_index := sort.Search(len(g.MapInfo.Bullets), func (i int) bool {
					return g.MapInfo.Bullets[i].GameObject.Id == target.GameObject.Id
				})
				var bullet = g.MapInfo.Bullets[bullet_index]
				// update the acceleration of two player
				var new_acceleration_a = AccelerationFormat {
					Up: 0.0,
					Down: 0.0,
					Left: 0.0,
					Right: 0.0,
				}
				new_acceleration_a.Up = player_session_a.Player.GameObject.Acceleration.Up * (-1.0)
				new_acceleration_a.Down = player_session_a.Player.GameObject.Acceleration.Down * (-1.0)
				new_acceleration_a.Left = player_session_a.Player.GameObject.Acceleration.Left * (-1.0)
				new_acceleration_a.Right = player_session_a.Player.GameObject.Acceleration.Right * (-1.0)

				player_session_a.Player.GameObject.Acceleration = new_acceleration_a
				
				// give the collision damage
				player_session_a.Player.Attr.HP -= float64(bullet.Damage) * 5.0
				// remove the bullet
				g.MapInfo.Bullets = append(g.MapInfo.Bullets[:bullet_index], g.MapInfo.Bullets[bullet_index+1:]...)
				break;
		}
	}
}

func GetObjectId (m GameObjectInterface) string {
	return m.GetId()
}


// define the JoinPlayer function in Game struct
func (g Game) JoinPlayer (session *PlayerSession) {
	// add the player session to channel
	g.JoinChannel <- session
}
