package main

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
	return Distance(p, Point{0, 0})
}

func Normalize(p Point) Point {
	return p.Div(int(Length(p)))
}
