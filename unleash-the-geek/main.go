package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

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

//FindRobots returns the player robots
func FindRobots(em EntityMap) []Entity {
	list := make([]Entity, 0)
	for _, e := range em {
		if e.EntityType == MyRobot && !e.Destroyed {
			list = append(list, e)
		}
	}
	return list
}

//ByDistance sorts Entities by Distance
func ByDistance(list []Entity, p Point) {
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
	tiles         TileMap
	entities      EntityMap
}

//NewGameObject creates GameObjects
func NewGameObject(input io.Reader) *GameObject {
	obj := &GameObject{}
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	obj.scanner = scanner

	scanner.Scan()
	fmt.Sscan(scanner.Text(), &obj.Width, &obj.Height)
	obj.tiles = make(TileMap)
	return obj
}

//ParseTurn collects the start inputs
func (o *GameObject) ParseTurn() {
	o.scanner.Scan()
	fmt.Sscan(o.scanner.Text(), &o.MyScore, &o.EnemyScore)
	for i := 0; i < o.Height; i++ {
		o.scanner.Scan()
		inputs := strings.Split(o.scanner.Text(), " ")
		line := ""
		for j := 0; j < o.Width; j++ {
			ore, err := strconv.ParseInt(inputs[2*j], 10, 32)
			line += inputs[2*j]
			if err != nil {
				ore = -1
			}
			hole, _ := strconv.ParseInt(inputs[2*j+1], 10, 32)
			point := Point{j, i}
			o.tiles[point] = Tile{point, ore, hole}
		}
		Debug("%v\n", line)
	}
	var entityCount int
	o.scanner.Scan()
	fmt.Sscan(o.scanner.Text(), &entityCount, &o.RadarCooldown, &o.TrapCooldown)
	o.entities = make(EntityMap)
	for i := 0; i < entityCount; i++ {
		var id, category, x, y, item int
		o.scanner.Scan()
		fmt.Sscan(o.scanner.Text(), &id, &category, &x, &y, &item)
		destroyed := x == -1 && y == -1
		o.entities[id] = Entity{Point{x, y}, id, category, destroyed, item}
	}
}

//TakeTurn returns the turns actions
func (o GameObject) TakeTurn() []string {
	actions := make([]string, 0)
	robots := FindRobots(o.entities)
	ByDistance(robots, Point{0, 0})
	for _, r := range robots {
		Debug("%#v -> %#v \n", r.Point, o.tiles[r.Add(Right)])
		actions = append(actions, "WAIT")
	}
	return actions
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
