package game

// import (
// 	"fmt"
// 	"time"
// 	"encoding/json"
// 	"github.com/syndtr/goleveldb/leveldb"
// 	"sync"
// )

type Point struct {
	X, Y float64
}
type Size struct {
	W, H float64
}

type BulletAttribute struct {
	Rotation int
	Size int
	Speed int
	Damage int
}
type Bullet struct {
	Attr BulletAttribute
	Position Point
	Owner string
}

type TrapAttribute struct {
	BulletSpeed int
	BulletDamage int
	BulletReload int
	BodyDamage int
}

type Trap struct {
	Position Point
	Rotation int
	Size int
	Attr TrapAttribute
}

type StuffAttribute struct {
	HP float64
	EXP int
	BodyDamage float64
}
type Stuff struct {
	Type int
	Position Point
	Attr StuffAttribute
}
