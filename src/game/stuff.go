package game

type StuffAttribute struct {
	HP float64
	EXP int
	BodyDamage float64
	
}
type Stuff struct {
	GameObject
	Type int
	Attr StuffAttribute
}