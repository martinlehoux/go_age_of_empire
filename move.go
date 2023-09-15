package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

type Move struct {
	IsActive    bool
	Origin      Point
	Destination Point
	Path        []Point
}

type Node struct {
	CostFromOrigin    int
	CostToDestination int
	Parent            Point
}

func (n Node) Cost() int {
	return n.CostFromOrigin + n.CostToDestination
}

type MoveMap struct {
	Width   int
	Height  int
	Blocked map[Point]bool
}

type PathSearch struct {
	origin      Point
	destination Point
	openList    map[Point]Node
	closedList  map[Point]Node
	moveMap     MoveMap
}

func SearchPath(origin Point, destination Point, moveMap MoveMap) ([]Point, bool) {
	search := PathSearch{
		origin:      origin,
		destination: destination,
		openList:    map[Point]Node{},
		closedList:  map[Point]Node{},
		moveMap:     moveMap,
	}
	return search.search()
}

func (s *PathSearch) addNeighbours(point Point, node Node) {
	if point.X >= 100 {
		s.consider(Point{point.X - 100, point.Y}, point)
	}
	if point.X <= s.moveMap.Width-100 {
		s.consider(Point{point.X + 100, point.Y}, point)
	}
	if point.Y >= 100 {
		s.consider(Point{point.X, point.Y - 100}, point)
	}
	if point.Y <= s.moveMap.Height-100 {
		s.consider(Point{point.X, point.Y + 100}, point)
	}
}

func (s *PathSearch) consider(neighbor Point, point Point) {
	if s.moveMap.Blocked[neighbor] {
		return
	}
	if _, ok := s.closedList[neighbor]; !ok {
		node := Node{
			CostFromOrigin:    s.closedList[point].CostFromOrigin + int(Distance(point, neighbor)),
			CostToDestination: int(Distance(neighbor, s.destination)),
			Parent:            point,
		}
		existing, ok := s.openList[neighbor]
		if !ok || existing.Cost() > node.Cost() {
			s.openList[neighbor] = node
		}
	}
}

// openList must not be empty
func (s *PathSearch) bestOpenNode() (Point, Node) {
	bestPoint := Point{}
	bestNode := Node{}
	bestCost := math.MaxInt32
	for point, node := range s.openList {
		if cost := node.Cost(); cost < bestCost {
			bestPoint = point
			bestNode = node
			bestCost = cost
		}
	}
	return bestPoint, bestNode
}

func (s *PathSearch) buildPath() []Point {
	point := s.destination
	path := []Point{s.destination}
	node := s.closedList[s.destination]
	for point != s.origin {
		point = node.Parent
		path = append(path, point)
		node = s.closedList[point]
	}
	slices.Reverse(path)
	return path
}

func (s *PathSearch) search() ([]Point, bool) {
	s.openList[s.origin] = Node{CostFromOrigin: 0, CostToDestination: int(Distance(s.origin, s.destination)), Parent: s.origin}
	point := s.origin
	var node Node
	for point != s.destination && len(s.openList) > 0 {
		point, node = s.bestOpenNode()
		s.closedList[point] = node
		delete(s.openList, point)
		s.addNeighbours(point, node)
	}
	if point == s.destination {
		return s.buildPath(), true
	}
	return []Point{}, false
}

func NewMove(origin Point, destination Point, moveMap MoveMap) Move {
	current := Point{
		origin.X / 100 * 100,
		origin.Y / 100 * 100,
	}
	search, ok := SearchPath(current, destination, moveMap)
	if !ok {
		slog.Info("no path found", slog.String("destination", destination.String()))
		return Move{IsActive: false}
	}
	return Move{IsActive: true, Origin: origin, Destination: destination, Path: search}
}

func (m *Move) Update(position Point, speed int) Point {
	if !m.IsActive {
		return position
	}
	next := m.Path[0]
	if Distance(position, next) < float64(speed) {
		if next == m.Destination {
			m.IsActive = false
			slog.Info("move finished", slog.String("destintation", m.Destination.String()))
		} else {
			m.Path = m.Path[1:]
		}
		return next
	}
	remainingMove := next.Sub(position)
	delta := remainingMove.Mul(speed).Div(int(Length(remainingMove)))
	return position.Add(delta)
}

func (e *Entity) StartMove(destination Point, moveMap MoveMap) {
	if e.Move.IsEnabled && e.Position.IsEnabled && e.Selection.IsEnabled {
		if e.Selection.Value.IsSelected {
			e.Move.Value = NewMove(e.Position.Value, destination, moveMap)
			slog.Info("entity starting move", slog.String("destination", destination.String()))
		}
	}
}

func (e *Entity) UpdateMove() {
	if e.Move.IsEnabled && e.Position.IsEnabled {
		e.Position.Value = e.Move.Value.Update(e.Position.Value, 10)
	}
}

func DrawMove(screen *ebiten.Image, e *Entity) {
	if e.Move.IsEnabled && e.Position.IsEnabled && e.Selection.IsEnabled {
		if e.Selection.Value.IsSelected && e.Move.Value.IsActive {
			last := e.Position.Value
			dx := +e.Image.Value.Bounds().Dx() / 2
			dy := +e.Image.Value.Bounds().Dy() / 2
			for _, point := range e.Move.Value.Path {
				vector.StrokeLine(screen, float32(last.X+dx), float32(last.Y+dy), float32(point.X+dx), float32(point.Y+dy), 10.0, color.RGBA{256 * 3 / 16, 256 * 3 / 16, 256 * 3 / 16, 256 / 4}, true)
				last = point
			}
		}
	}
}
