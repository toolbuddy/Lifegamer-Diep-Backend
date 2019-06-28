package game


import (
	"math/rand"
	"io/ioutil"
	"os"
	"strconv"
	"encoding/json"
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
)

/**
 * Diep:
 * The struct of player diep.
 *
 * @property {string} Name					 												- the name of the player
 * @property {GameObject} 					 												- the game object struct of the diep
 */
 type Diep struct {
	Name string
	GameObject
}

/**
 * Bullet:
 * The struct of player bullet.
 *
 * @property {GameObject} 					 												- the game object struct of the diep
 * @property {int} Damage																		- the damage of the bullet
 * @property {int} Existence																- the existence time of the bullet
 * @property {string} Owner					 												- the name of the owner
 */
type Bullet struct {
	GameObject
	Damage int
	Existence int
	Owner string
}

/**
 * StuffAttribute:
 * The struct of stuff attribute.
 *
 * @property {float64} HP					 													- the HP of the stuff
 * @property {int} EXP																			- the EXP of the stuff
 * @property {float64} BodyDamage														- the body damage of the stuff
 */
type StuffAttribute struct {
	HP float64
	EXP int
	BodyDamage float64
	
}

/**
 * Stuff:
 * The struct of stuff.
 *
 * @property {GameObject} 					 												- the game object struct of the stuff
 * @property {int} Type																			- the type number of the stuff
 * @property {StuffAttrbute} Attr														- the attribute of the stuff
 */
type Stuff struct {
	GameObject
	Type int
	Attr StuffAttribute
}

/**
 * TrapAttribute:
 * The struct of trap attribute.
 *
 * @property {int} BulletSpeed 					 										- the bullet speed of the trap
 * @property {int} BulletDamage															- the bullet damage of the trap
 * @property {int} BulletReload															- the bullet reload time of the trap
 * @property {int} BodyDamage																- the body damage of the trap
 */
type TrapAttribute struct {
	BulletSpeed int
	BulletDamage int
	BulletReload int
	BodyDamage int
}

/**
 * Trap:
 * The struct of trap.
 *
 * @property {GameObject} 					 												- the game object struct of the trap
 * @property {TrapAttribute} Attr														- the attribute of the trap
 */
type Trap struct {
	GameObject
	Attr TrapAttribute
}

/**
 * CollisionDetection:
 * The struct to keep the reference in a collision.
 *
 * @property {*Diep} object_a					 											- the reference of Diep
 * @property {GameObjectInterface} object_b					 				- the interface of collider
 */
type CollisionDetection struct {
	object_a *Diep
	object_b GameObjectInterface
}

/**
 * Map:
 * The struct of game map to keep all item reference.
 *
 * @property {[]*Diep} Dieps					 												- the slice of Dieps on the field
 * @property {[]*Bullet} Bullets					 										- the slice of Bullets on the field
 * @property {[]*Stuff} Stuffs					 											- the slice of Stuffs on the field
 * @property {[]*Trap} Traps					 												- the slice of Traps on the field
 * @property {[]CollisionDetection} Collisions					 			- the slice of CollisionDetection on the field
 */
type Map struct {
	Dieps []*Diep
	Bullets []*Bullet
	Stuffs []*Stuff
	Traps []*Trap
	Collisions []CollisionDetection
}

/**
 * <game>.NewStuff:
 * The function to new a stuff in random position.
 *
 * @return {*Stuff}
 */
func NewStuff() *Stuff {
	// read the stuff type through json file
	jsonFile, err := os.Open("src/config/stuffType.json")
	if (err != nil) {
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// decode the json to the map[string]interface{}
	var stuffType map[string]interface{}
  json.Unmarshal([]byte(byteValue), &stuffType)

	// generate the random type of the stuff
	var type_num = int(rand.Float64() * 4 + 1)
	var type_attr = stuffType[strconv.Itoa(type_num)].(StuffAttribute)

	uuid, _ := util.NewUUID()
	var new_stuff = Stuff {
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
		Type: type_num,
		Attr: StuffAttribute {
			HP: type_attr.HP,
			EXP: type_attr.EXP,
			BodyDamage: type_attr.BodyDamage,
		},
	}
	return &new_stuff
}

// TODO: event driven to generate the stuff