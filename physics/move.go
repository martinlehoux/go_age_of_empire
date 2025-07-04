package physics

import (
	"age_of_empires/ecs"
	"fmt"
	"math"

	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

const MOVE_SPEED = 10.0

type Path []Point

func (p Path) String() string {
	return fmt.Sprintf("%s -> %s", p[0].String(), p[len(p)-1].String())
}

type Move struct {
	IsActive    bool
	Origin      Point
	Destination Point
	Path        Path
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

func SearchPath(origin Point, destination Point, moveMap MoveMap) (Path, bool) {
	search := PathSearch{
		origin:      origin,
		destination: destination,
		openList:    map[Point]Node{},
		closedList:  map[Point]Node{},
		moveMap:     moveMap,
	}
	return search.search()
}

func (s *PathSearch) addNeighbors(point Point, node Node) {
	if point.X >= 100 {
		s.consider(Point{X: point.X - 100, Y: point.Y}, point)
	}
	if point.X <= s.moveMap.Width-100 {
		s.consider(Point{X: point.X + 100, Y: point.Y}, point)
	}
	if point.Y >= 100 {
		s.consider(Point{X: point.X, Y: point.Y - 100}, point)
	}
	if point.Y <= s.moveMap.Height-100 {
		s.consider(Point{X: point.X, Y: point.Y + 100}, point)
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
	bestPoint := Point{X: 0, Y: 0}
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

func (s *PathSearch) buildPath() Path {
	point := s.destination
	path := Path{s.destination}
	node := s.closedList[s.destination]
	for point != s.origin {
		point = node.Parent
		path = append(path, point)
		node = s.closedList[point]
	}
	slices.Reverse(path)
	return path
}

func (s *PathSearch) search() (Path, bool) {
	s.openList[s.origin] = Node{CostFromOrigin: 0, CostToDestination: int(Distance(s.origin, s.destination)), Parent: s.origin}
	point := s.origin
	var node Node
	for point != s.destination && len(s.openList) > 0 {
		point, node = s.bestOpenNode()
		s.closedList[point] = node
		delete(s.openList, point)
		s.addNeighbors(point, node)
	}
	if point == s.destination {
		return s.buildPath(), true
	}
	return Path{}, false
}

func NewMove(origin Point, destination Point, moveMap MoveMap) Move {
	current := Point{
		X: origin.X / 100 * 100,
		Y: origin.Y / 100 * 100,
	}
	path, ok := SearchPath(current, destination, moveMap)
	if !ok {
		slog.Info("no path found", slog.String("destination", destination.String()))
		return Move{IsActive: false}
	}
	return Move{IsActive: true, Origin: origin, Destination: destination, Path: path[1:]}
}

func (p *Path) isValid(moveMap MoveMap) bool {
	for _, point := range *p {
		if moveMap.Blocked[point] {
			return false
		}
	}
	return true
}

func StartMove(move *ecs.Component[Move], position ecs.Component[Point], destination Point, moveMap MoveMap) {
	if move.IsEnabled && position.IsEnabled {
		move.Value = NewMove(position.Value, destination, moveMap)
		slog.Info("entity starting move", slog.String("destination", destination.String()))
	}
}

func moveToward(position *ecs.Component[Point], destination Point) {
	remainingMove := destination.Sub(position.Value)
	delta := remainingMove.Mul(MOVE_SPEED).Div(int(Length(remainingMove)))
	position.Value = position.Value.Add(delta)
}

func UpdateMove(move *ecs.Component[Move], position *ecs.Component[Point], moveMap MoveMap) {
	if !move.IsEnabled || !position.IsEnabled || !move.Value.IsActive {
		return
	}
	if len(move.Value.Path) == 0 {
		move.Value.IsActive = false
		return
	}
	next := move.Value.Path[0]
	if Distance(position.Value, next) < MOVE_SPEED {
		position.Value = next
		remainingPath := move.Value.Path[1:]
		if next == move.Value.Destination {
			move.Value.IsActive = false
			slog.Info("move finished", slog.String("destination", move.Value.Destination.String()))
		} else if remainingPath.isValid(moveMap) {
			move.Value.Path = remainingPath
		} else {
			// Here maybe the destination has to change
			path, ok := SearchPath(next, move.Value.Destination, moveMap)
			if !ok {
				slog.Info("no path found", slog.String("destination", move.Value.Destination.String()))
				move.Value.IsActive = false
			} else {
				move.Value.Path = path[1:]
			}
		}
	} else {
		moveToward(position, next)
	}
}
