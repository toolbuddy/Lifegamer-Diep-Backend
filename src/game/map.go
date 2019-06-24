package game


type CollisionDetection struct {
	object_a *Diep
	object_b GameObjectInterface
}

type Map struct {
	Dieps []*Diep
	Bullets []*Bullet
	Stuffs []*Stuff
	Traps []*Trap
	Collisions []CollisionDetection
}
