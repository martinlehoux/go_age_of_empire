package main

type Vec struct {
	X, Y float64
}

type Rectangle struct {
	leftTop     Vec
	rightBottom Vec
}

func (r Rectangle) Contains(point Vec) bool {
	return point.X > r.leftTop.X && point.X < r.rightBottom.X && point.Y > r.leftTop.Y && point.Y < r.rightBottom.Y
}
