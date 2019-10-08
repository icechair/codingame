package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	Width  = 12
	Height = 12
	HQ     = 0
	Me     = 0
	Enemy  = 1
	NoOne  = -1
)

type Point struct {
	X int
	Y int
}
type Entity interface {
	Position() Point
}
type Mine struct {
	Point
}

func NewMine() *Mine {
	var x, y int
	fmt.Scan(&x, &y)
	return &Mine{Point{x, y}}
}

type Tile struct {
	Point
	Owner      int
	active     bool
	MineSpot   *Mine
	OccupiedBy Entity
}

func NewTile(x, y int, tile string, mines map[Point]*Mine) *Tile {
	if tile == "#" {
		return nil
	}
	owner := NoOne
	if strings.ToLower(tile) == "o" {
		owner = Me
	} else if strings.ToLower(tile) == "x" {
		owner = Enemy
	}
	active := tile == strings.ToUpper(tile)

	t := &Tile{Point{x, y}, owner, active, nil, nil}

	if mine, ok := mines[t.Point]; ok {
		t.MineSpot = mine
	}
	return t
}

type Building struct {
	Point
	Owner        int
	BuildingType int
}

func (b Building) Position() Point {
	return b.Point
}

func NewBuilding() *Building {
	var owner, buildingType, x, y int
	fmt.Scan(&owner, &buildingType, &x, &y)
	return &Building{Point{x, y}, owner, buildingType}
}

type Unit struct {
	Point
	Owner int
	ID    int
	Level int
}

func (u Unit) Position() Point {
	return u.Point
}

func NewUnit() *Unit {
	var owner, unitID, level, x, y int
	fmt.Scan(&owner, &unitID, &level, &x, &y)
	return &Unit{Point{x, y}, owner, unitID, level}
}

type TileMap map[Point]*Tile

type GameState struct {
	Gold        int
	Income      int
	EnemyGold   int
	EnemyIncome int
	Map         TileMap
	MyBuildings map[int]*Building
	MyUnits     map[int]*Unit
	EnemyBuildings
}

func NewGameState(mines map[Point]*Mine) *GameState {
	gs := &GameState{}

	fmt.Scan(&gs.Gold)
	fmt.Scan(&gs.Income)

	fmt.Scan(&gs.EnemyGold)
	fmt.Scan(&gs.EnemyIncome)
	gs.Map = make(TileMap)
	for i := 0; i < 12; i++ {
		var line string
		fmt.Scan(&line)
		cols := strings.Split(line, "")
		for j, c := range cols {
			tile := NewTile(j, i, c, mines)
			if tile != nil {
				gs.Map[tile.Point] = tile
			}
		}
	}
	var buildingCount int
	fmt.Scan(&buildingCount)

	for i := 0; i < buildingCount; i++ {
		building := NewBuilding()
		gs.Buildings = append(gs.Buildings, building)
		gs.Map[building.Point].OccupiedBy = building
	}
	var unitCount int
	fmt.Scan(&unitCount)

	for i := 0; i < unitCount; i++ {
		unit := NewUnit()
		gs.Units = append(gs.Units, unit)
		gs.Map[unit.Point].OccupiedBy = unit
	}
	return gs
}

func (p Point) distance(b Point) int {
	d := p.X - b.X + p.Y - b.Y
	if d < 0 {
		return d * -1
	}
	return d
}
func (m TileMap) neighbours(p Point) []*Tile {
	var tiles []*Tile
	if t, ok := m[Point{p.X + 1, p.Y}]; ok {
		tiles = append(tiles, t)
	}
	if t, ok := m[Point{p.X - 1, p.Y}]; ok {
		tiles = append(tiles, t)
	}
	if t, ok := m[Point{p.X, p.Y + 1}]; ok {
		tiles = append(tiles, t)
	}
	if t, ok := m[Point{p.X, p.Y - 1}]; ok {
		tiles = append(tiles, t)
	}
	return tiles
}

func (m TileMap) direction(from, to Point) *Tile {
	if from == to {
		return nil
	}
	neighbours := m.neighbours(from)
	sort.SliceStable(neighbours, func(a, b int) bool {
		return to.distance(neighbours[a].Point) < to.distance(neighbours[b].Point)
	})
	for _, n := range neighbours {
		if n.OccupiedBy == nil {
			return n
		}
	}
	return neighbours[0]
}

func (gs GameState) turn() []string {
	var actions []string

	actions = append(actions, "WAIT")
	const hqNeigbours = gs.Map.neighbours()
	for _, tile := range hqNeigbours {

	}
	return actions
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var numberMineSpots int
	fmt.Scan(&numberMineSpots)
	var mines map[Point]*Mine
	for i := 0; i < numberMineSpots; i++ {
		mine := NewMine()
		mines[mine.Point] = mine
	}

	for {
		gs := NewGameState(mines)
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println(strings.Join(gs.turn(), ";")) // Write action to stdout
	}
}
