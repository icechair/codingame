package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
TODOS:
 - Optimize the tile storage. only interesting tiles are those with stuff in or on it
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
func (a Point) Clamp(o GameObject) Point {
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

//Distance Returns the manhatten distance between Points
func (a Point) Distance(b Point) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
}

//Game Constants
const (
	Hole       = 1
	Empty      = 0
	MyRobot    = 0
	EnemyRobot = 1
	MyRadar    = 2
	MyTrap     = 3
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

//Tile holds the game map information
type Tile struct {
	Point
	Ore  int64
	Hole int64
}

//TileMap map of Tiles by Point
type TileMap = map[Point]Tile

//Entity are Game Elements(robots, items)
type Entity struct {
	Point
	ID         int
	EntityType int
	Destroyed  bool
	Item       int
}

//EntityMap map of Entities By ID
type EntityMap = map[int]Entity

//FindMyRobots returns the player robots
func FindMyRobots(em EntityMap) []Entity {
	list := make([]Entity, 0)
	for _, e := range em {
		if e.EntityType == MyRobot {
			list = append(list, e)
		}
	}
	return list
}

//EntitiesByDistance sorts Entities by Distance
func EntitiesByDistance(list []Entity, p Point) {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Distance(p) < list[j].Distance(p)
	})
}

//EntitiesByID sorts Entities by ID
func EntitiesByID(list []Entity) {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
}

//FindOreTiles returns ore tiles
func FindOreTiles(tm TileMap) []Tile {
	list := make([]Tile, 0)
	for _, t := range tm {
		if t.Ore > 0 {
			list = append(list, t)
		}
	}
	return list
}

//TilesByDistance sorts Tiles by Distance
func TilesByDistance(list []Tile, p Point) {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Distance(p) < list[j].Distance(p)
	})
}

//GameObject holds the game information
type GameObject struct {
	scanner       *bufio.Scanner
	Width         int
	Height        int
	MyScore       int
	EnemyScore    int
	RadarCooldown int
	TrapCooldown  int
	Tiles         TileMap
	Entities      EntityMap
	TileHash      uint64
	RadarHolder   int
	TrapHolder    int
}

//NewGameObject creates GameObjects
func NewGameObject(input io.Reader) *GameObject {
	obj := &GameObject{}
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	obj.scanner = scanner

	scanner.Scan()
	fmt.Sscan(scanner.Text(), &obj.Width, &obj.Height)
	obj.Tiles = make(TileMap)
	return obj
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

//ParseTurn collects the start inputs
func (o *GameObject) ParseTurn() {
	o.scanner.Scan()
	fmt.Sscan(o.scanner.Text(), &o.MyScore, &o.EnemyScore)
	tilestring := ""
	oreline := ""
	holeline := ""

	for i := 0; i < o.Height; i++ {
		o.scanner.Scan()
		line := o.scanner.Text()
		inputs := strings.Split(line, " ")
		//tilestring += line
		for j := 0; j < o.Width; j++ {
			ore, err := strconv.ParseInt(inputs[2*j], 10, 32)
			oreline += inputs[2*j]
			if err != nil {
				ore = -1
			}
			hole, _ := strconv.ParseInt(inputs[2*j+1], 10, 32)
			holeline += inputs[2*j+1]
			point := Point{j, i}
			o.Tiles[point] = Tile{point, ore, hole}
		}
		oreline += "\n"
		tilestring += oreline
		holeline += "\n"

	}
	o.TileHash = hash(tilestring)
	//Debug("%#v\n", o.TileHash)
	//Debug("ORE:\n%v\n", oreline)
	//Debug("HOLE:\n%v\n", holeline)
	var entityCount int
	o.scanner.Scan()
	fmt.Sscan(o.scanner.Text(), &entityCount, &o.RadarCooldown, &o.TrapCooldown)
	o.Entities = make(EntityMap)
	o.RadarHolder = Nothing
	o.TrapHolder = Nothing
	for i := 0; i < entityCount; i++ {
		var id, category, x, y, item int
		o.scanner.Scan()
		fmt.Sscan(o.scanner.Text(), &id, &category, &x, &y, &item)
		destroyed := x == -1 && y == -1
		entity := Entity{Point{x, y}, id, category, destroyed, item}
		o.Entities[id] = entity
		if category == MyRobot {
			if item == Radar {
				o.RadarHolder = id
			}
			if item == Trap {
				o.TrapHolder = id
			}
		}
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

//TakeTurn returns the turns actions
func (o *GameObject) TakeTurn() []string {
	actions := make([]string, 0)
	robots := FindMyRobots(o.Entities)
	oreTiles := FindOreTiles(o.Tiles)
	if len(oreTiles) > 0 {
		rand.Seed(int64(len(oreTiles) * len(robots) * 123123))

	} else {
		rand.Seed(int64(o.TileHash))
	}
	a := Point{random(0, o.Width/2), random(0, o.Height/2)}
	b := Point{random(0, o.Width/2), random(0, o.Height/2)}
	center := a.Add(b)
	EntitiesByID(robots)

	for _, robot := range robots {
		TilesByDistance(oreTiles, robot.Point)
		Debug("%#v -> %#v : %#v\n", robot, center, robot.Distance(center))
		action := o.RobotAction(robot, robots, oreTiles, center)
		if strings.Contains(action, "DIG") && len(oreTiles) > 0 {
			oreTiles = oreTiles[1:]
		}
		actions = append(actions, action)
	}
	return actions
}

//RobotAction decides the action a robot takes
func (o *GameObject) RobotAction(robot Entity, robots []Entity, oreTiles []Tile, p Point) string {
	// robot := o.Entities[id]
	if robot.Destroyed {
		return "WAIT"
	}
	if o.RadarHolder == Nothing && len(oreTiles) < len(robots)*2 {
		o.RadarHolder = robot.ID
		return "REQUEST RADAR"
	}
	if robot.Item == Ore {
		return fmt.Sprintf("MOVE 0 %v ID(%v)", robot.Y, robot.ID)
	}
	if o.RadarHolder == robot.ID {
		return fmt.Sprintf("DIG %v %v ID(%v)", p.X, p.Y, robot.ID)
	}
	if len(oreTiles) > 0 {
		tile := oreTiles[0]
		//if tile.Hole == 0 {
		return fmt.Sprintf("DIG %v %v ID(%v)", tile.X, tile.Y, robot.ID)
		//}
	}
	offset := Point{random(-RadarRange, RadarRange), random(-RadarRange, RadarRange)}
	target := p.Add(offset).Clamp(*o)

	return fmt.Sprintf("MOVE %v %v ID(%v)", target.X, target.Y, robot.ID)
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
	Up = Point{0, -1}
	Down = Point{0, 1}
	Left = Point{-1, 0}
	Right = Point{1, 0}

	game := NewGameObject(os.Stdin)
	for {
		game.ParseTurn()
		actions := game.TakeTurn()
		for _, action := range actions {
			fmt.Println(action) // WAIT|MOVE x y|DIG x y|REQUEST item
		}
	}
}
