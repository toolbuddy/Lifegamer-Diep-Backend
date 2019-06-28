package game

import (
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
)

/**
 * GameObject:
 * The struct of game object.
 *
 * @property {string} Id					 														- the unique identity between game objects
 * @property {util.Point} Position														- the position on screen
 * @property {float64} Mass																		- the mass of the game object
 * @property {float64} Radius																	- the collision circle area radius
 * @property {util.VelocityFormat} Velocity										- the velocity of the game object
 * @property {util.AccelerationFormat} Acceleration						- the velocity of the game object
 */
 type GameObject struct {
	Id string
	Position util.Point
	Mass float64
	Radius float64
	Velocity util.VelocityFormat
	Acceleration util.AccelerationFormat
}

/**
 * GameObjectInterface:
 * The interface of game object.
 *
 * @function {string} GetId					 													- the function to get the game object id
 */
type GameObjectInterface interface {
	GetId() string
}

/**
 * <*GameObject>.GetId:
 * The function to get the game object id
 * 
 * @return {string}
 */
func (g *GameObject) GetId() string {
	return g.Id
}

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
 */
type Game struct {
	Name string
	Sessions []*PlayerSession
	MapInfo Map
	JoinChannel chan *PlayerSession
	Field *util.Size
	Framerate float64
}