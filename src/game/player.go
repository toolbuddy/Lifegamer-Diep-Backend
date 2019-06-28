package game

import (
	"time"
	"math/rand"
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
)

/**
 * PlayerAttribute:
 * The struct of player attribute.
 *
 * @property {string} Name					 													- the name of the player
 * @property {time.Time} CreatedAt														- the join time of the player
 * @property {int} Score																			- the total score of the player
 * @property {int} Level																			- the level of the player
 * @property {int} EXP																				- the current EXP of the player
 * @property {float64} HP																			- the current HP of the player
 */
type PlayerAttribute struct {
	Name string
	CreatedAt time.Time
	Score int
	Level int
	EXP int
	HP float64
}

/**
 * PlayerStatus:
 * The struct of player status.
 *
 * @property {int} MaxHP					 														- the max HP level of the player
 * @property {int} HPRegeneration															- the HP eegeneration level of the player
 * @property {int} MoveSpeed																	- the move speed level of the player
 * @property {int} BulletSpeed																- the bullet speed level of the player
 * @property {int} BulletPenetration													- the bullet penetration level of the player
 * @property {int} BulletReload																- the bullet reload level of the player
 * @property {int} BulletDamage																- the bullet damage level of the player
 * @property {int} BodyDamage																	- the body damage level of the player
 */
type PlayerStatus struct {
	MaxHP int
	HPRegeneration int
	MoveSpeed int
	BulletSpeed int
	BulletPenetration int
	BulletReload int
	BulletDamage int
	BodyDamage int
}

/**
 * Player:
 * The struct of player.
 *
 * @property {GameObject} 					 													- the game object struct of the player
 * @property {PlayerAttribute} Attr														- the struct of the player attribute
 * @property {PlayerStatus} Status														- the struct of the player status
 */
type Player struct {
	GameObject
	Attr PlayerAttribute
	Status PlayerStatus
}
