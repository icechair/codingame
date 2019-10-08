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
	Destroyed  bool
	Item       int
}

//TileMap represents a map of interesting Points on the Playfield
type TileMap = map[Point]int

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
	MyRobots      []Entity
	EnemyRobots   []Entity
}

//OreTiles returns a slice of Points
func (ti TurnInput) OreTiles() []Point {
	list := make([]Point, 0)
	for p, ore := range ti.RadarTiles {
		if ore > 0 {
			list = append(list, p)
		}
	}
	return list
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
			1000000*len(ti.OreTiles()),
		//10000000*len(ti.HoleTiles),
	)
}

//NewTurnInput creates new TurnInputs out of scanner
func NewTurnInput(scanner *bufio.Scanner, game *GameObject) TurnInput {
	o := TurnInput{}
	o.RadarHolder = Nothing
	o.TrapHolder = Nothing
	o.RadarTiles = make(TileMap)
	o.HoleTiles = make(TileMap)
	o.MyRobots = make([]Entity, 0)
	o.EnemyRobots = make([]Entity, 0)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &o.MyScore, &o.EnemyScore)
	oreline := ""
	holeline := ""
	for i := 0; i < game.Height; i++ {
		scanner.Scan()
		line := scanner.Text()
		inputs := strings.Split(line, " ")
		//tilestring += line
		for j := 0; j < game.Width; j++ {
			point := Point{j, i}

			ore, err := strconv.ParseInt(inputs[2*j], 10, 32)
			oreline += inputs[2*j]
			if err == nil {
				o.RadarTiles[point] = int(ore)
			}
			hole, _ := strconv.ParseInt(inputs[2*j+1], 10, 32)
			holeline += inputs[2*j+1]
			if hole > 0 {
				o.HoleTiles[point] = int(hole)
			}
			//o.Tiles[point] = Tile{point, int(ore), int(hole)}
		}
		oreline += "\n"
		holeline += "\n"

	}
	//Debug("%#v\n", o.TileHash)
	//Debug("ORE:\n%v\n", oreline)
	//Debug("HOLE:\n%v\n", holeline)
	var entityCount int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &entityCount, &o.RadarCooldown, &o.TrapCooldown)

	for i := 0; i < entityCount; i++ {
		var id, category, x, y, item int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &id, &category, &x, &y, &item)
		destroyed := x == -1 && y == -1
		point := Point{x, y}
		entity := Entity{point, id, category, destroyed, item}
		if category == MyRobot {
			if item == Radar {
				o.RadarHolder = id
			}
			if item == Trap {
				o.TrapHolder = id
			}
			o.MyRobots = append(o.MyRobots, entity)
		} else if category == EnemyRobot {
			o.EnemyRobots = append(o.EnemyRobots, entity)
		} else if category == MyRadar {
			o.HoleTiles[point] = MyRadar
		} else if category == MyTrap {
			o.HoleTiles[point] = MyTrap
		}
	}
	return o
}

//PointsByDistance sorts a slice of points by distance to a point
func PointsByDistance(list []Point, p Point) {
	sort.SliceStable(list, (func(i, j int) bool {
		return p.Distance(list[i]) < p.Distance(list[j])
	}))
}

//TakeTurn returns the turns actions
func (o *GameObject) TakeTurn(ti *TurnInput) []string {
	actions := make([]string, 0)
	rand.Seed(ti.ToSeed())
	a := Point{random(0, o.Width/2), random(0, o.Height/2)}
	b := Point{random(0, o.Width/2), random(0, o.Height/2)}
	center := a.Add(b)
	oreTiles := ti.OreTiles()

	for _, r := range ti.MyRobots {
		PointsByDistance(oreTiles, r.Point)
		action := o.RobotAction(r, ti, oreTiles, center)
		if strings.Contains(action, "DIG") && len(oreTiles) > 0 {
			oreTiles = oreTiles[1:]
		}
		actions = append(actions, action)
	}
	return actions
}

//RobotAction decides the action a robot takes
func (o *GameObject) RobotAction(robot Entity, ti *TurnInput, oreTiles []Point, p Point) string {
	// robot := o.Entities[id]
	if robot.Destroyed {
		return "WAIT"
	}
	if ti.RadarHolder == Nothing && ti.RadarCooldown == 0 {
		ti.RadarHolder = robot.ID
		return "REQUEST RADAR"
	}
	if robot.Item == Ore {
		return fmt.Sprintf("MOVE 0 %v ID(%v)", robot.Y, robot.ID)
	}
	if ti.RadarHolder == robot.ID {
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
