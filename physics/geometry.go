package physics

import (
	"image"
	"math"
)

type Point = image.Point
type Rectangle = image.Rectangle

func Distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2))
}

func Length(p Point) float64 {
	return Distance(p, Point{X: 0, Y: 0})
}

func Normalize(p Point) Point {
	return p.Div(int(Length(p)))
}

func AdjacentPoints(from Point) []Point {
	return []Point{
		from.Add(Point{X: 0, Y: -100}),
		from.Add(Point{X: 0, Y: 100}),
		from.Add(Point{X: -100, Y: 0}),
		from.Add(Point{X: 100, Y: 0}),
	}
}

func Closest(from Point, candidates []Point, moveMap MoveMap) (Point, int) {
	var closest Point
	closestDistance := math.MaxInt
	for _, dest := range candidates {
		path, ok := SearchPath(from, dest, moveMap)
		if !ok {
			continue
		}
		if len(path) < closestDistance {
			closest = dest
			closestDistance = len(path)
		}
	}
	return closest, closestDistance
}
