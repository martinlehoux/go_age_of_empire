package main

type Order interface {
	Update(e *Entity, g *Game)
}

func (e *Entity) UpdateOrder(g *Game) {
	if !e.Order.IsEnabled || e.Order.Value == nil {
		return
	}
	e.Order.Value.Update(e, g)
}
