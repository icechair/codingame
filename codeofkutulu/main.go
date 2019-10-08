package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var width int
var height int

func debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

//Const stuff
const (
	WALL     = rune('#')
	SPAWN    = rune('w')
	EMPTY    = rune('.')
	EXPLORER = 0
	WANDERER = 1
)

//Point stuff
type Point struct {
	x, y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

//Field stuff
type Field [][]rune

//NewField stuff
func NewField(row, col int) Field {
	r := make(Field, row)
	for k := range r {
		r[k] = make([]rune, col)
	}
	return r
}

func (f Field) String() string {
	out := ""
	for _, r := range f {
		for _, c := range r {
			out = out + fmt.Sprintf("%s", string(c))
		}
		out = out + "\n"
	}
	return out
}

//Neighbours Stuff
func (f Field) Neighbours(p Point) []Point {
	list := []Point{
		Point{p.x - 1, p.y},
		Point{p.x + 1, p.y},
		Point{p.x, p.y - 1},
		Point{p.x, p.y + 1},
	}
	return list
}

//DistanceMatrix stuff
type DistanceMatrix map[string]int

//NewDistanceMatrix stuff
func NewDistanceMatrix(p Point, f Field) DistanceMatrix {
	dm := make(DistanceMatrix)

	return dm
}
func (dm DistanceMatrix) String() string {
	out := ""
	for k, v := range dm {
		out += fmt.Sprintf("'%s':%d\n", k, v)
	}
	return out
}

//Entity stuff
type Entity struct {
	Point
	hash       string
	entityType int
	id         int
	param0     int // sanity, time to spawn, time to recall
	param1     int // ??, minion state
	param2     int // ??, target
}

//NewEntity makes
func NewEntity(entityType string, id, col, row, param0, param1, param2 int) *Entity {
	etype := EXPLORER
	if entityType == "WANDERER" {
		etype = WANDERER
	}
	hash := fmt.Sprintf("%d-%d", etype, id)
	return &Entity{Point{col, row}, hash, etype, id, param0, param1, param2}
}

func (e Entity) String() string {
	return fmt.Sprintf("E(%s, %d, %d, [%d,%d], %d, %d, %d)", e.hash, e.entityType, e.id, e.x, e.y, e.param0, e.param1, e.param2)
}

//Entities stuff
type Entities map[string]*Entity

func (el Entities) filter(fn func(e Entity) bool) []*Entity {
	out := make([]*Entity, 0)
	for _, e := range el {
		if fn(*e) {
			out = append(out, e)
		}
	}
	return out
}

//Turn stuff
func (e Entity) Turn(field Field, entities Entities) string {
	explorers := entities.filter(
		func(t Entity) bool {
			return t.entityType == EXPLORER && t.hash != e.hash
		},
	)
	enemies := entities.filter(
		func(t Entity) bool {
			return t.entityType != EXPLORER
		},
	)
	debug("%#v, %#v\n", len(explorers), len(enemies))
	return "WAIT"
}

func manhattan(a, b Point) int {
	dx := math.Abs(float64(b.x - a.x))
	dy := math.Abs(float64(b.y - a.y))
	return int(dx + dy)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	fmt.Sscan(scanner.Text(), &width)

	scanner.Scan()
	fmt.Sscan(scanner.Text(), &height)
	field := NewField(height, width)
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		runes := []rune(line)
		field[i] = runes
	}

	// sanityLossLonely: how much sanity you lose every turn when alone, always 3 until wood 1
	// sanityLossGroup: how much sanity you lose every turn when near another player, always 1 until wood 1
	// wandererSpawnTime: how many turns the wanderer take to spawn, always 3 until wood 1
	// wandererLifeTime: how many turns the wanderer is on map after spawning, always 40 until wood 1
	var sanityLossLonely, sanityLossGroup, wandererSpawnTime, wandererLifeTime int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &sanityLossLonely, &sanityLossGroup, &wandererSpawnTime, &wandererLifeTime)

	for {
		// entityCount: the first given entity corresponds to your explorer
		var entityCount int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &entityCount)
		entities := make(Entities)
		for i := 0; i < entityCount; i++ {
			var entityType string
			var id, x, y, param0, param1, param2 int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &entityType, &id, &x, &y, &param0, &param1, &param2)
			entity := NewEntity(entityType, id, x, y, param0, param1, param2)
			entities[entity.hash] = entity
			debug("%s\n", entity)
		}
		myExplorer := entities["0-0"]

		fmt.Println(myExplorer.Turn(field, entities)) // MOVE <x> <y> | WAIT
	}
}
