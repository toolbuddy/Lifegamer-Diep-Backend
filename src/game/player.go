package game

import (
	"time"
	"math/rand"
	"github.com/f26401004/Lifegamer-Diep-backend/src/utils"
)

type PlayerAttribute struct {
	Name string
	CreatedAt time.Time
	Score int
	Level int
	EXP int
	HP float64
}
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
type Player struct {
	GameObject
	Attr PlayerAttribute
	Status PlayerStatus
}

// define the NewPlayer function in game package
func NewPlayer (name string) *Player {
	uuid, _ := utils.NewUUID()
	var new_player = Player {
		GameObject: GameObject {
			Id: uuid,
			Position: Point {
				X: rand.Float64() * 1023,
				Y: rand.Float64() * 1023,
			},
			Mass: 1.0,
			Radius: 50.0,
			Velocity: VelocityFormat {
				X: 0.0,
				Y: 0.0,
			},
			Acceleration: AccelerationFormat {
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