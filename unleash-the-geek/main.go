package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

/*
TODOS:
 - Optimize the radar placement. keep radar positions, radar coverage in mind, find out if ore placement is somehow biased
 - optimize dig target aquisition. dont consider empty holes, dig randomly if no ore is available, keep move distance in mind
 - lookup enemy radars and traps. tile history, item history for robots
*/

//Debug prints messages
func Debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

//Point is a tile coordinate
type Point struct {
	X int
	Y int
}

//Add 2 points
func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

//Clamp restricts the point into the boundaries of the gameworld
func (a Point) Clamp(o *GameObject) Point {
	if a.X < 0 {
		a.X = 0
	}
	if a.Y < 0 {
		a.Y = 0
	}
	if a.X >= o.Width {
		a.X = o.Width - 1
	}
	if a.Y >= o.Height {
		a.Y = o.Height - 1
	}
	return a
}

//Distance returns the manhatten distance between points
func (a Point) Distance(b Point) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
}

//Neigbours returns the neighbouring cell tiles
func (a Point) Neigbours(o *GameObject) []Point {
	list := make([]Point, 4)
	list[0] = a.Add(Up).Clamp(o)
	list[1] = a.Add(Down).Clamp(o)
	list[2] = a.Add(Left).Clamp(o)
	list[3] = a.Add(Right).Clamp(o)

	return list
}

//Game Constants
const (
	Hole       = 1
	Empty      = 0
	MyRobot    = 0
	EnemyRobot = 1
	MyRadar    = 2
	MyTrap     = 3
	EnemyTrap  = 8
	EnemyRadar = 9
	Nothing    = -1
	Radar      = 2
	Trap       = 3
	Ore        = 4
	MoveRange  = 4
	RadarRange = 4
)

//Up Direction
var Up Point

//Down Direction
var Down Point

//Left Direction
var Left Point

//Right Direction
var Right Point

func init() {
	Up = Point{0, -1}
	Down = Point{0, 1}
	Left = Point{-1, 0}
	Right = Point{1, 0}
}

//Entity are Game Elements(robots, items)
type Entity struct {
	Point
	ID         int
	EntityType int
	Item       int
}

//Destroyed robots have position{-1,-1}
func (e Entity) Destroyed() bool {
	return e.X == Nothing && e.Y == Nothing
}

//Idle returns true if the robot is not carrying an item
func (e Entity) Idle() bool {
	return e.Item == Nothing
}

//TileMap represents a map of interesting Points on the Playfield
type TileMap map[Point]int

//EntityMap storage
type EntityMap map[int]Entity

//GameObject holds the game information
type GameObject struct {
	Width   int
	Height  int
	History []TurnInput
}

//NewGameObject creates GameObjects
func NewGameObject(scanner *bufio.Scanner) GameObject {
	obj := GameObject{}
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &obj.Width, &obj.Height)
	obj.History = make([]TurnInput, 0)
	return obj
}

//TurnInput defines the Input each Turn
type TurnInput struct {
	MyScore       int
	EnemyScore    int
	RadarCooldown int
	TrapCooldown  int
	RadarHolder   int
	TrapHolder    int
	RadarTiles    TileMap
	HoleTiles     TileMap
	MyRobots      EntityMap
	EnemyRobots   EntityMap
}

//OreTiles returns a slice of Points
func (ti TurnInput) OreTiles() []Point {
	list := make([]Point, 0)
	for p, ore := range ti.RadarTiles {
		if ore > 0 {
			list = append(list, p)
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		return ti.RadarTiles[list[i]] > ti.RadarTiles[list[j]]
	})
	return list
}

//PlayerRobots gets robots sorted by id
func (ti TurnInput) PlayerRobots() []Entity {
	var keys []int
	for k := range ti.MyRobots {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	robots := make([]Entity, 0)
	for _, id := range keys {
		robots = append(robots, ti.MyRobots[id])
	}
	return robots
}

//ToSeed returns the turn hash
func (ti TurnInput) ToSeed() int64 {
	return int64(
		ti.MyScore +
			10*ti.EnemyScore +
			100*ti.RadarCooldown +
			1000*ti.TrapCooldown +
			10000*ti.RadarHolder +
			100000*ti.TrapHolder +
			1000000*len(ti.OreTiles()) +
			10000000*len(ti.HoleTiles),
	)
}

//NewTurnInput creates new TurnInputs out of scanner
func NewTurnInput(scanner *bufio.Scanner, game *GameObject) TurnInput {
	ti := TurnInput{}
	ti.RadarHolder = Nothing
	ti.TrapHolder = Nothing
	ti.RadarTiles = make(TileMap)
	ti.HoleTiles = make(TileMap)
	ti.MyRobots = make(EntityMap)
	ti.EnemyRobots = make(EntityMap)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &ti.MyScore, &ti.EnemyScore)

	oreline := ""
	holeline := ""
	for i := 0; i < game.Height; i++ {
		scanner.Scan()
		line := scanner.Text()
		//Debug("%v\n", line)
		inputs := strings.Split(line, " ")

		for j := 0; j < game.Width; j++ {
			point := Point{j, i}

			ore, err := strconv.ParseInt(inputs[2*j], 10, 32)
			oreline += inputs[2*j]
			if err == nil {
				ti.RadarTiles[point] = int(ore)
			}
			hole, _ := strconv.ParseInt(inputs[2*j+1], 10, 32)
			holeline += inputs[2*j+1]
			if hole > 0 {
				ti.HoleTiles[point] = int(hole)
			}
			//o.Tiles[point] = Tile{point, int(ore), int(hole)}
		}
		oreline += "\n"
		holeline += "\n"
	}
	//Debug("ORE:\n%v\nHOLE:\n%v\n", oreline, holeline)
	var entityCount int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &entityCount, &ti.RadarCooldown, &ti.TrapCooldown)

	for i := 0; i < entityCount; i++ {
		var id, category, x, y, item int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &id, &category, &x, &y, &item)
		point := Point{x, y}
		entity := Entity{point, id, category, item}
		if category == MyRobot {
			if item == Radar {
				ti.RadarHolder = id
			}
			if item == Trap {
				ti.TrapHolder = id
			}
			ti.MyRobots[id] = entity
		} else if category == EnemyRobot {
			ti.EnemyRobots[id] = entity
		} else if category == MyRadar {
			ti.HoleTiles[point] = MyRadar
		} else if category == MyTrap {
			ti.HoleTiles[point] = MyTrap
		}
	}
	/* you cant see what the enemy robots are carrying
	if len(game.History) > 0 {
		lt := game.History[len(game.History)-1]
		for p, hole := range lt.HoleTiles {
			if hole == EnemyTrap {
				o.HoleTiles[p] = EnemyTrap
			}
		}

		for id, robot := range o.EnemyRobots {
			if robot.Item == Nothing {
				if lt.EnemyRobots[id].Item == Trap {
					for _, p := range robot.Neigbours(game) {
						if _, ok := o.HoleTiles[p]; ok {
							o.HoleTiles[p] = EnemyTrap
						}
					}
				}

				if lt.EnemyRobots[id].Item == Radar {
					for _, p := range robot.Neigbours(game) {
						if _, ok := o.HoleTiles[p]; ok {
							o.HoleTiles[p] = EnemyRadar
						}
					}
				}
			}
		}
	}
	*/
	/* for p, h := range o.HoleTiles {
		Debug("%v, %v\n", p, h)
	}
	for _, r := range o.EnemyRobots {
		Debug("%#v\n", r)
	} */
	return ti
}

//PointsByDistance sorts a slice of points by distance to a point
func PointsByDistance(list []Point, p Point) {
	sort.SliceStable(list, (func(i, j int) bool {
		return p.Distance(list[i]) < p.Distance(list[j])
	}))
}

//TakeTurn returns the turns actions
func (o GameObject) TakeTurn(ti *TurnInput) []string {
	actions := make([]string, 0)
	rand.Seed(42)
	oreTiles := ti.OreTiles()

	for _, r := range ti.PlayerRobots() {
		PointsByDistance(oreTiles, r.Point)
		action := o.RobotAction(r, ti, oreTiles)
		if strings.Contains(action, "DIG") && len(oreTiles) > 0 {
			oreTiles = oreTiles[1:]
		}
		actions = append(actions, action)
	}
	return actions
}

//RobotAction decides the action a robot takes
func (o GameObject) RobotAction(robot Entity, ti *TurnInput, oreTiles []Point) string {
	if robot.Destroyed() {
		return "WAIT"
	}
	if robot.Item == Ore {
		return fmt.Sprintf("MOVE 0 %v ID(%v)", robot.Y, robot.ID)
	}
	if robot.Item == Radar {
		return o.RadarAction(robot, ti)
	}
	if robot.Item == Trap {
		return o.TrapAction(robot, ti)
	}
	if ti.RadarHolder == Nothing && ti.RadarCooldown == 0 && robot.Point.Y > 0 && len(oreTiles) < 2*len(ti.MyRobots) {
		ti.RadarHolder = robot.ID
		return "REQUEST RADAR"
	}
	if ti.TrapHolder == Nothing && ti.TrapCooldown == 0 {
		ti.TrapHolder = robot.ID
		return "REQUEST TRAP"
	}
	if len(oreTiles) > 0 {
		var tile Point
		for len(oreTiles) > 0 {
			tile, oreTiles = oreTiles[0], oreTiles[1:]
			if hole, ok := ti.HoleTiles[tile]; !ok || (ok && hole != MyTrap) {
				break
			}
			//if tile.Hole == 0 {
		}
		return fmt.Sprintf("DIG %v %v ID(%v)", tile.X, tile.Y, robot.ID)
		//}
	}
	target := Point{1, robot.Y}
	for {
		target = target.Add(Point{random(1, MoveRange), 0})
		if _, ok := ti.HoleTiles[target]; !ok {
			break
		}
	}

	return fmt.Sprintf("DIG %v %v ID(%v)", target.X, target.Y, robot.ID)
}

//TrapAction places Traps
func (o GameObject) TrapAction(robot Entity, ti *TurnInput) string {
	targets := make([]Point, 0)
	for k := range ti.RadarTiles {
		targets = append(targets, k)
	}
	PointsByDistance(targets, robot.Point)
	target := robot.Add(
		Point{
			random(-1, 1),
			random(-1, 1),
		},
	)
	for len(targets) > 0 {
		target, targets = targets[0], targets[1:]
		if ti.RadarTiles[target] > 1 {
			break
		}
	}

	return fmt.Sprintf("DIG %v %v ID(%v)", target.X, target.Y, robot.ID)
}

//RadarAction places Radars
func (o GameObject) RadarAction(robot Entity, ti *TurnInput) string {
	p := Point{5, robot.Y}
	targets := make([]Point, 0)
	for p, hole := range ti.HoleTiles {
		if hole == MyRadar {
			rd := p.Add(Point{4, 5}) //.Clamp(&o)
			if _, ok := ti.HoleTiles[rd]; !ok && rd.X < o.Width && rd.Y < o.Height {
				targets = append(targets, rd.Clamp(&o))
			}
			ru := p.Add(Point{4, -5}) //.Clamp(&o)
			if _, ok := ti.HoleTiles[ru]; !ok && ru.X < o.Width && ru.Y > 0 {
				targets = append(targets, ru.Clamp(&o))
			}
		}
	}
	if len(targets) > 0 {
		PointsByDistance(targets, robot.Point)
		p = targets[0]
	}
	return fmt.Sprintf("DIG %v %v ID(%v)", p.X, p.Y, robot.ID)
}

/*
Action order for one turn

1. If DIG commands would trigger Traps, they go off.
2. The other DIG commands are resolved.
3. REQUEST commands are resolved.
4. Request timers are decremented.
5. MOVE and WAIT commands are resolved.
6. Ore is delivered to the headquarters.
*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	game := NewGameObject(scanner)
	for {
		ti := NewTurnInput(scanner, &game)
		Debug("RadarLength: %v, HoleLength: %v\n", len(ti.RadarTiles), len(ti.HoleTiles))
		actions := game.TakeTurn(&ti)
		game.History = append(game.History, ti)
		for _, action := range actions {
			fmt.Println(action) // WAIT|MOVE x y|DIG x y|REQUEST item
		}

	}
}
