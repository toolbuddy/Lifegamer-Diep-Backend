package util

/**
 * Point:
 * The struct to present the point.
 *
 * @property {float64} X 									- the x of the point
 * @property {float64} Y									- the y of the point
 */
type Point struct {
	X, Y float64
}

/**
 * Size:
 * The struct to present the size.
 *
 * @property {float64} W 									- the width of the size
 * @property {float64} H									- the height of the size
 */
type Size struct {
	W, H float64
}

/**
 * AccelerationFormat:
 * The struct to present the acceleration.
 *
 * @property {float64} Up 								- the up volume of the acceleration
 * @property {float64} Down								- the down volume of the acceleration
 * @property {float64} Left								- the left volume of the acceleration
 * @property {float64} Right							- the right volume of the acceleration
 */
type AccelerationFormat struct {
	Up float64
	Down float64
	Left float64
	Right float64
}

/**
 * VelocityFormat:
 * The struct to present the velocity.
 *
 * @property {float64} X 								- the X volume of the velocity
 * @property {float64} Y								- the Y volume of the velocity
 */
type VelocityFormat struct {
	X float64
	Y float64
} 

/**
 * MoveDirection:
 * The struct to keep the player current moving direction.
 *
 * @property {bool} Up 									- the move up status of the player
 * @property {bool} Down								- the move down status of the player
 * @property {bool} Left								- the move left status of the player
 * @property {bool} Right								- the move right status of the player
 */
 type MoveDirection struct {
	Up bool
	Down bool
	Left bool
	Right bool
}
