package main

type Move struct {
	IsActive    bool
	Destination Point
}

func (m *Move) Update(position Point, speed int) Point {
	if !m.IsActive {
		return position
	}
	remainingMove := m.Destination.Sub(position)
	delta := remainingMove.Mul(speed).Div(int(Length(remainingMove)))
	next := position.Add(delta)
	if Distance(next, m.Destination) < float64(speed) {
		m.IsActive = false
	}
	return next
}
