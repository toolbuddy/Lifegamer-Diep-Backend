package game


import (
	"math/rand"
	"math"
	"time"
	"io/ioutil"
	"os"
	"strconv"
	"encoding/json"
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
	"log"
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
	HP int
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
 * <*Game>.NewStuff:
 * The function to new a stuff in random position.
 *
 * @return {*Stuff}
 */
func (g *Game) NewStuff() *Stuff {
	// read the stuff type through json file
	jsonFile, err := os.Open("./src/config/stuffType.json")
	if (err != nil) {
		log.Print(err)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// decode the json to the map[string]interface{}
	var stuffType map[string]StuffAttribute
	json.Unmarshal([]byte(byteValue), &stuffType)

	// generate the random type of the stuff
	var type_num = int(rand.Float64() * 5 + 1)
	var type_attr = stuffType[strconv.Itoa(type_num)]

	uuid, _ := util.NewUUID()
	var new_stuff = Stuff {
		GameObject: GameObject {
			Id: uuid,
			Position: util.Point {
				X: rand.Float64() * g.Field.W,
				Y: rand.Float64() * g.Field.H,
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

func (g *Game) generateStuff () {
	for i := 0; i < 50 ; i++ {
		target := g.NewStuff()
		if (target == nil) {
			return
		}
		g.MapInfo.Stuffs = append(g.MapInfo.Stuffs, target)
		// log.Printf("Generate Stuff %d at Position X: %f, Y: %f, Total stuff: %d\n", target.Type, target.GameObject.Position.X, target.GameObject.Position.Y, len(g.MapInfo.Stuffs))
	}
	for {
		// sleep correspond to the number of the stuffs
		var stepDelay = int32(math.Pow(float64(len(g.MapInfo.Stuffs)), 1.618 / math.Max(float64(len(g.Sessions)), 2) * math.Max(float64(len(g.Sessions) - 1), 1)))
		time.Sleep(time.Duration(stepDelay) * time.Millisecond)
		// if the stuff meet the maximum number, the stop generate the stuff
		if (float64(len(g.MapInfo.Stuffs)) > math.Max(float64(len(g.Sessions)), 1) * 1000 * 1.618 / math.Max(float64(len(g.Sessions)), 2) * math.Max(float64(len(g.Sessions) - 1), 1)) {
			return
		}
		// generate the stuff and append to the map
		target := g.NewStuff()
		if (target == nil) {
			return
		}
		g.MapInfo.Stuffs = append(g.MapInfo.Stuffs, target)
		// log.Printf("Generate Stuff %d at Position X: %f, Y: %f, Total stuff: %d\n", target.Type, target.GameObject.Position.X, target.GameObject.Position.Y, len(g.MapInfo.Stuffs))
	}
}