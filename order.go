package main

type Order interface {
	Update(e *Entity, g *Game)
}

func (e *Entity) UpdateOrder(g *Game) {
	if e.Order.IsEnabled {
		if e.Order.Value != nil {
			e.Order.Value.Update(e, g)
		}
	}
}
