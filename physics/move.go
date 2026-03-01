package physics

import (
	"fmt"
	"math"

	"age_of_empires/ecs"

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

var neighborsDirectionsDiagonal = []Point{
	{X: -100, Y: 0},    // left
	{X: 100, Y: 0},     // right
	{X: 0, Y: -100},    // up
	{X: 0, Y: 100},     // down
	{X: -100, Y: -100}, // top-left
	{X: 100, Y: -100},  // top-right
	{X: -100, Y: 100},  // bottom-left
	{X: 100, Y: 100},   // bottom-right
}

// Add neighbors to the open list
// Choose neighborsDirections
func (s *PathSearch) addNeighbors(point Point) {
	for _, dir := range neighborsDirectionsDiagonal {
		neighbor := Point{X: point.X + dir.X, Y: point.Y + dir.Y}
		if neighbor.X >= 0 && neighbor.X <= s.moveMap.Width-100 &&
			neighbor.Y >= 0 && neighbor.Y <= s.moveMap.Height-100 {
			s.consider(neighbor, point)
		}
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
		s.addNeighbors(point)
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

func moveToward(position Point, destination Point) Point {
	remainingMove := destination.Sub(position)
	delta := remainingMove.Mul(MOVE_SPEED).Div(int(Length(remainingMove)))
	return position.Add(delta)
}

func UpdateMoveSystem(moveMap MoveMap, masks []ecs.Mask, moves []Move, positions []Point) {
	required := ecs.CM_Position | ecs.CM_Move
	for i, mask := range masks {
		if mask&required != required {
			continue
		}
		move := &moves[i]
		position := &positions[i]
		if !move.IsActive {
			continue
		}
		if len(move.Path) == 0 {
			move.IsActive = false
			continue
		}
		next := move.Path[0]
		if Distance(*position, next) >= MOVE_SPEED {
			*position = moveToward(*position, next)
			continue
		}
		// About to snap to next cell — check it's still free
		if moveMap.Blocked[next] {
			slog.Info("cell blocked on arrival, repathing", slog.String("next", next.String()))
			path, ok := SearchPath(*position, move.Destination, moveMap)
			if !ok {
				slog.Info("no path found", slog.String("destination", move.Destination.String()))
				move.IsActive = false
				continue
			}
			move.Path = path[1:]
			continue
		}
		oldCell := Point{X: position.X / 100 * 100, Y: position.Y / 100 * 100}
		*position = next
		moveMap.Blocked[next] = true
		delete(moveMap.Blocked, oldCell)
		remainingPath := move.Path[1:]
		if next == move.Destination {
			move.IsActive = false
			slog.Info("move finished", slog.String("destination", move.Destination.String()))
			continue
		}
		if remainingPath.isValid(moveMap) {
			move.Path = remainingPath
			continue
		}
		// Here maybe the destination has to change
		path, ok := SearchPath(next, move.Destination, moveMap)
		if !ok {
			slog.Info("no path found", slog.String("destination", move.Destination.String()))
			move.IsActive = false
			continue
		}
		move.Path = path[1:]
	}
}
