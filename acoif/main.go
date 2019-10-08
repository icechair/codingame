package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

//Constants
const (
	DEBUG      = true
	Width      = 12
	Height     = 12
	HQ         = 0
	MINE       = 1
	TOWER      = 2
	ME         = 0
	ENEMY      = 1
	UNOCCUPIED = -1
	UNIT       = 0
	BUILDING   = 1
)

//UnitCost struct
type UnitCost struct {
	Level  int
	Train  int
	Upkeep int
}

var dirs = []Point{
	P(1, 0),
	P(0, 1),
	P(-1, 0),
	P(0, -1),
}

var unitPrices = make(map[int]UnitCost)

// fmt.Fprintln(os.Stderr, "Debug messages...")
func debug(format string, a ...interface{}) {
	if DEBUG {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func init() {
	unitPrices[1] = UnitCost{1, 10, 1}
	unitPrices[2] = UnitCost{2, 20, 4}
	unitPrices[3] = UnitCost{3, 30, 20}

}

func main() {
	var numberMineSpots int
	fmt.Scan(&numberMineSpots)
	mines := make(map[Point]Point)
	for i := 0; i < numberMineSpots; i++ {
		mine := ParseMineSpot()
		mines[mine] = mine
	}
	n := 0
	for {
		gs := ParseGameState(mines)
		fmt.Println(strings.Join(gs.Turn(n), ";")) // Write action to stdout
		n++
	}
}

//Point Structure
type Point struct {
	X int
	Y int
}

//P returns a new Point
func P(x, y int) Point {
	return Point{x, y}
}
func (p Point) add(b Point) Point {
	return Point{p.X + b.X, p.Y + b.Y}
}

func (p Point) distance(b Point) int {
	d := p.X - b.X + p.Y - b.Y
	if d < 0 {
		return d * -1
	}
	return d
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// Entity Game Object interface
type Entity interface {
	Position() Point
}

//ParseMineSpot from STDIN
func ParseMineSpot() Point {
	var x, y int
	fmt.Scan(&x, &y)
	return Point{x, y}
}

//Tile Gamemap structure
type Tile struct {
	Point
	Owner      int
	active     bool
	MineSpot   bool
	OccupiedBy Entity
	DefendedBy *Building
}

func (t Tile) String() string {
	return fmt.Sprintf("T(%s, %d, %v)", t.Point, t.Owner, t.active)
}

//NewTile constructor
func NewTile(x, y int, tile string, mines map[Point]Point) *Tile {
	if tile == "#" {
		return nil
	}
	owner := UNOCCUPIED
	if strings.ToLower(tile) == "o" {
		owner = ME
	} else if strings.ToLower(tile) == "x" {
		owner = ENEMY
	}
	active := tile == strings.ToUpper(tile)

	t := &Tile{Point{x, y}, owner, active, false, nil, nil}

	if _, ok := mines[t.Point]; ok {
		t.MineSpot = true
	}
	return t
}

//Building Entity
type Building struct {
	Point
	Owner        int
	BuildingType int
}

//Position Entity implementation
func (b Building) Position() Point {
	return b.Point
}

//NewBuilding Constructor
func NewBuilding() *Building {
	var owner, buildingType, x, y int
	fmt.Scan(&owner, &buildingType, &x, &y)
	return &Building{Point{x, y}, owner, buildingType}
}

//Unit Entity
type Unit struct {
	Point
	Owner int
	ID    int
	Level int
}

//Position Entity implementation
func (u Unit) Position() Point {
	return u.Point
}

//ParseUnit from STDIN
func ParseUnit() *Unit {
	var owner, unitID, level, x, y int
	fmt.Scan(&owner, &unitID, &level, &x, &y)
	return &Unit{Point{x, y}, owner, unitID, level}
}

//TileMap game playfield
type TileMap map[Point]*Tile

func (m TileMap) neighbours(p Point) []*Tile {
	var tiles []*Tile
	for _, dir := range dirs {
		neighbor := p.add(dir)
		if tile, ok := m[neighbor]; ok {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

//BFS breath first search
func (m *TileMap) BFS(start Point, until func(c Point) bool) map[Point]int {
	frontier := []Point{start}
	distance := make(map[Point]int)
	distance[start] = 0
	for len(frontier) > 0 {
		var current Point
		current, frontier = frontier[0], frontier[1:]
		if until(current) {
			return distance
		}
		for _, next := range m.neighbours(current) {
			if _, ok := distance[next.Point]; !ok {
				frontier = append(frontier, next.Point)
				distance[next.Point] = 1 + distance[current]
			}
		}
	}
	return distance
}

//ByDistance Sorter
type ByDistance struct {
	distance map[Point]int
	tiles    []*Tile
}

func (b ByDistance) Len() int      { return len(b.tiles) }
func (b ByDistance) Swap(i, j int) { b.tiles[i], b.tiles[j] = b.tiles[j], b.tiles[i] }
func (b ByDistance) Less(i, j int) bool {
	if b.distance[b.tiles[i].Point] < b.distance[b.tiles[j].Point] {
		return true
	}
	if b.tiles[i].Y < b.tiles[j].Y {
		return true
	}
	return b.tiles[i].X < b.tiles[j].X
}

//TilesSortedByDistanceFrom point
func (m TileMap) TilesSortedByDistanceFrom(p Point) []*Tile {
	distance := m.BFS(p, func(c Point) bool { return false })
	tiles := make([]*Tile, len(m))
	idx := 0
	for _, tile := range m {
		tiles[idx] = tile
		idx++
	}
	sort.Stable(ByDistance{distance, tiles})
	return tiles
}

//BuildingTypeMap hashmap
type BuildingTypeMap = map[int][]*Building

//GameState structure
type GameState struct {
	Gold           int
	Income         int
	EnemyGold      int
	EnemyIncome    int
	Map            TileMap
	Buildings      BuildingTypeMap
	Units          []*Unit
	EnemyBuildings BuildingTypeMap
	EnemyUnits     []*Unit
}

var enemyPos Point

//ParseGameState from STDIN
func ParseGameState(mines map[Point]Point) *GameState {
	gs := &GameState{}
	fmt.Scan(&gs.Gold)
	fmt.Scan(&gs.Income)

	fmt.Scan(&gs.EnemyGold)
	fmt.Scan(&gs.EnemyIncome)
	gs.Map = make(TileMap)
	gs.Buildings = make(BuildingTypeMap)
	gs.EnemyBuildings = make(BuildingTypeMap)
	gs.Units = make([]*Unit, 0)
	gs.EnemyUnits = make([]*Unit, 0)
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
		var buildList map[int][]*Building
		buildList = gs.Buildings
		if building.Owner == ENEMY {
			buildList = gs.EnemyBuildings
		}
		if _, ok := buildList[building.BuildingType]; !ok {
			buildList[building.BuildingType] = make([]*Building, 0)
		}
		buildList[building.BuildingType] = append(buildList[building.BuildingType], building)
		if building.BuildingType == TOWER {
			next := gs.Map.neighbours(building.Point)
			for _, tile := range next {
				if tile.Owner == building.Owner && tile.active {
					tile.DefendedBy = building
				}
			}
		}
		gs.Map[building.Point].OccupiedBy = *building
	}
	var unitCount int
	fmt.Scan(&unitCount)

	for i := 0; i < unitCount; i++ {
		unit := ParseUnit()
		if unit.Owner == ME {
			gs.Units = append(gs.Units, unit)

		} else if unit.Owner == ENEMY {
			gs.EnemyUnits = append(gs.EnemyUnits, unit)
		}
		gs.Map[unit.Point].OccupiedBy = *unit
	}

	enemyPos = gs.EnemyBuildings[HQ][0].Point
	return gs
}

func (g GameState) mineCost() int {
	return 20 + 4*len(g.Buildings[MINE])
}

//Train action
func (g *GameState) Train(tile *Tile) string {
	var level = 1
	if tile.DefendedBy != nil {
		level = 3
	} else if u, ok := tile.OccupiedBy.(Unit); ok {
		level = u.Level + 1
	} else if b, ok := tile.OccupiedBy.(Building); ok {
		if b.BuildingType == TOWER {
			level = 3
		}
	}
	if level > 3 {
		level = 3
	}
	if g.Gold >= unitPrices[level].Train && g.Income >= unitPrices[level].Upkeep {
		unit := &Unit{tile.Point, ME, -1, level}
		g.Gold -= unitPrices[level].Train
		g.Income -= unitPrices[level].Upkeep
		if tile.Owner != ME || tile.Owner == ME && !tile.active {
			g.Income++
		}
		tile.OccupiedBy = *unit
		g.Units = append(g.Units, unit)
		debug("TRAINING: %#v\n", tile)
		return fmt.Sprintf("TRAIN %d %d %d", level, tile.X, tile.Y)
	}
	return ""
}

//Move action
func (g *GameState) Move(unit *Unit, enemyHQ Point) string {
	if unit.ID == -1 {
		return ""
	}
	var target *Tile
	// debug("unit:%#v\n", unit)
	until := func(c Point) bool {
		if c == unit.Point {
			return false
		}
		target = g.Map[c]
		// debug("\ttarget: %#v\n", target)
		if target.OccupiedBy == nil {
			if target.Owner != ME || (target.Owner == ME && !target.active) {
				return true
			}
			return false
		}
		if u, ok := target.OccupiedBy.(Unit); ok {
			// debug("\tbyUnit:%#v\n", target.OccupiedBy)
			if unit.Level == 3 {
				return true
			}
			if u.Level < unit.Level {
				return true
			}
			return false
		}
		if b, ok := target.OccupiedBy.(Building); ok {
			if b.Owner == ME {
				return false
			}
			if b.BuildingType != TOWER {
				return true
			}
			if b.BuildingType == TOWER && unit.Level == 3 {
				return true
			}
			return false
		}
		if target.DefendedBy == nil || target.DefendedBy != nil && unit.Level == 3 {
			return true
		}
		return false
	}
	g.Map.BFS(unit.Point, until)

	if until(target.Point) {
		g.Map[unit.Point].OccupiedBy = nil
		g.Map[target.Point].OccupiedBy = *unit
		g.Map[target.Point].Owner = ME
		g.Map[target.Point].active = true
		return fmt.Sprintf("MOVE %d %d %d", unit.ID, target.X, target.Y)
	}
	return ""
}

//Build action
func (g *GameState) Build(tile *Tile, myHQ Point) string {
	if tile.MineSpot {
		if tile.OccupiedBy == nil {
			if g.Gold > g.mineCost() {
				g.Gold -= g.mineCost()
				g.Income += 4
				building := &Building{tile.Point, ME, MINE}
				g.Map[tile.Point].OccupiedBy = building
				return fmt.Sprintf("BUILD MINE %d %d", tile.X, tile.Y)
			}
		}
	}
	if tile.distance(myHQ) == 2 {
		if tile.OccupiedBy == nil && !tile.MineSpot {
			if tile.Owner == ME {
				if g.Gold > 15 {
					//g.Gold -= 15
					//building := &Building{tile.Point, ME, TOWER}
					//g.Map[tile.Point].OccupiedBy = building
					//return fmt.Sprintf("BUILD TOWER %d %d", tile.X, tile.Y)
				}
			}
		}
	}
	return ""
}

func (g GameState) String() string {
	return fmt.Sprintf("Gold:%d, Income:%d, EnemyGold:%d,EnemyIncome: %d", g.Gold, g.Income, g.EnemyGold, g.EnemyIncome)
}

//Turn actions
func (g *GameState) Turn(turn int) []string {
	var actions []string
	debug("turn:%d, %s\n", turn, g)
	myHQ := g.Buildings[HQ][0]

	enemyHQ := g.EnemyBuildings[HQ][0]
	// tilesFromEnemy := g.Map.TilesSortedByDistanceFrom(enemyHQ.Point)

	for _, unit := range g.Units {
		action := g.Move(unit, enemyHQ.Point)
		if action != "" {
			actions = append(actions, action)
		}
	}
	tilesFromPlayer := g.Map.TilesSortedByDistanceFrom(myHQ.Point)
	for _, tile := range tilesFromPlayer {
		if tile.Owner == ME && tile.active {
			action := g.Build(tile, myHQ.Point)
			if action != "" {
				actions = append(actions, action)
			}
		}
		if tile.Owner != ME || (tile.Owner == ME && !tile.active) {
			for _, next := range g.Map.neighbours(tile.Point) {
				if next.Owner == ME && next.active {
					action := g.Train(tile)
					if action != "" {
						actions = append(actions, action)
					}
				}
			}
		}

	}

	//	debug("fromPlayer:\n%s\n", tilesFromPlayer)
	//	debug("\nfromEnemy:\n%s\n", tilesFromEnemy)
	actions = append(actions, "WAIT")
	debug("turn:%d, %s\n", turn, g)
	return actions
}
