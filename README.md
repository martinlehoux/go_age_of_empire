## TODO

```go
func (g \*Game) updatePlacingWall(cursor Point) {
  if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
    if cursor.In(soil.Bounds()) {
      position := cursor.Sub(soil.Position).Div(100).Mul(100).Add(Point{50, 50})
      g.Blocks = append(g.Blocks, Block{NewTile(Point{100, 100}, position, color.RGBA{0x00, 0x00, 0x00, 0xff}), WallBlock})
    }
  }
}
```

```go
  if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && cursor.In(wallButton.CollisionBounds()) {
    if g.CurrentAction == PlacingWall {
      g.CurrentAction = Selecting
    } else {
      g.CurrentAction = PlacingWall
    }
  }
```
