package game

import (
	"time"
	"math/rand"
)
type PlayerAttribute struct {
	Name string
	CreatedAt time.Time
	Score int
	Level int
	EXP int
	HP float64
	Position Point
	Angle int
	Velocity float64
	Acceleration float64
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
	Attr PlayerAttribute
	Status PlayerStatus
}

// define the NewPlayer function in game package
func NewPlayer (name string) *Player {
	var new_player = Player {
		Attr: PlayerAttribute {
			Name: name,
			CreatedAt: time.Now(),
			Score: 0,
			Level: 1,
			EXP: 0,
			HP: 100,
			Position: Point {
				X: rand.Float64() * 1023,
				Y: rand.Float64() * 1023,
			},
			Angle: 0,
			Velocity: 0.0,
			Acceleration: 0.0,
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