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
type AccelerationFormat struct {
	Up float64
	Down float64
	Left float64
	Right float64
}
type VelocityFormat struct {
	X float64
	Y float64
} 
type GameObject struct {
	Id string
	Position Point
	Mass float64
	Radius float64
	Velocity VelocityFormat
	Acceleration AccelerationFormat
}

type GameObjectInterface interface {
	GetId() string
}


func (g *GameObject) GetId() string {
	return g.Id
}