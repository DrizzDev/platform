package point

import "errors"

const limit = 1 << 16

// Point is a device-pixel coordinate. Neither axis alone is a location, so the pair is validated together.
type Point struct {
	x int
	y int
}

type Input struct {
	X int
	Y int
}

func New(input Input) (Point, error) {
	point := Point{x: input.X, y: input.Y}
	if failure := point.validate(); failure != nil {
		return Point{}, failure
	}
	return point, nil
}

func (point Point) X() int {
	return point.x
}

func (point Point) Y() int {
	return point.y
}

func (point Point) validate() error {
	switch {
	case point.x < 0 || point.y < 0:
		return errors.New("point coordinates must not be negative")
	case point.x > limit || point.y > limit:
		return errors.New("point coordinates exceed the maximum")
	default:
		return nil
	}
}
