package util

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

/**
 * MoveDirection:
 * The struct to keep the player current moving direction.
 *
 * @property {bool} Up 									- all dieps in player views
 * @property {bool} Down								- all stuffs in player views
 * @property {bool} Left								- all bullets in player views
 * @property {bool} Right								- all traps in player views
 */
 type MoveDirection struct {
	Up bool
	Down bool
	Left bool
	Right bool
}