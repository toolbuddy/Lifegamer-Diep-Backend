package game

type TrapAttribute struct {
	BulletSpeed int
	BulletDamage int
	BulletReload int
	BodyDamage int
}

type Trap struct {
	GameObject
	Rotation int
	Attr TrapAttribute
}