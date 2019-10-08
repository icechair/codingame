package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const MAX_W = 1920
const MAX_H = 750

func debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

type Point struct {
	x int
	y int
}

func arrivalTime(distance int, movementSpeed int) float64 {
	return float64(distance / movementSpeed)
}

func (p Point) String() string {
	return fmt.Sprintf("%d %d", p.x, p.y)
}

func (a Point) Distance(b Point) int {
	dx := float64(b.x - a.x)
	dy := float64(b.y - a.y)
	return int(math.Sqrt(dx*dx + dy*dy))
}

type Entity struct {
	entityType string
	pos        Point
	radius     int
}
type Entities []*Entity

type Item struct {
	name             string
	cost             int
	damage           int
	health           int
	maxHealth        int
	mana             int
	maxMana          int
	moveSpeed        int
	manaRegeneration int
	isPotion         bool
}
type Items []*Item

type Unit struct {
	id               int
	team             int
	unitType         string
	pos              Point
	attackRange      int
	health           int
	maxHealth        int
	shield           int
	attackDamage     int
	movementSpeed    int
	stunDuration     int
	goldValue        int
	countdown1       int
	countdown2       int
	countdown3       int
	mana             int
	maxMana          int
	manaRegeneration int
	heroType         string
	isVisible        bool
	itemsOwned       int
}
type Units []*Unit

func (l Units) First() *Unit {
	return l[0]
}
func (l Units) Heroes(team int) Units {
	heroes := make(Units, 0)
	for _, unit := range l {
		if unit.unitType == "HERO" && unit.team == team {
			heroes = append(heroes, unit)
		}
	}
	return heroes
}

type GameState struct {
	team      int
	entities  Entities
	items     Items
	gold      int
	enemyGold int
	roundType int
	units     Units
}

func Wait() string                            { return "WAIT" }
func Move(p Point) string                     { return fmt.Sprintf("MOVE %s", p) }
func Attack(target *Unit) string              { return fmt.Sprintf("ATTACK %d", target.id) }
func AttackNearest(unitType string) string    { return fmt.Sprintf("ATTACK_NEAREST %s", unitType) }
func MoveAttack(p Point, target *Unit) string { return fmt.Sprintf("MOVE_ATTACK %s %d", p, target.id) }
func Buy(item *Item) string                   { return fmt.Sprintf("BUY %s", item.name) }
func Sell(item *Item) string                  { return fmt.Sprintf("SELL %s", item.name) }

func NewGame() *GameState {
	var myTeam int
	fmt.Scan(&myTeam)
	var bushAndSpawnPointCount int
	fmt.Scan(&bushAndSpawnPointCount)
	entities := make(Entities, bushAndSpawnPointCount)
	for i := 0; i < bushAndSpawnPointCount; i++ {
		var entityType string
		var x, y, radius int
		fmt.Scan(&entityType, &x, &y, &radius)
		entities[i] = &Entity{entityType, Point{x, y}, radius}
	}
	var itemCount int
	fmt.Scan(&itemCount)
	items := make(Items, itemCount)
	for i := 0; i < itemCount; i++ {
		var itemName string
		var itemCost, damage, health, maxHealth, mana, maxMana, moveSpeed, manaRegeneration int
		var isPotion bool
		fmt.Scan(&itemName, &itemCost, &damage, &health, &maxHealth, &mana, &maxMana, &moveSpeed, &manaRegeneration, &isPotion)
		items[i] = &Item{itemName, itemCost, damage, health, maxHealth, mana, maxMana, moveSpeed, manaRegeneration, isPotion}
	}
	return &GameState{
		myTeam,
		entities,
		items,
		0,
		0,
		-1,
		make(Units, 0),
	}
}

func (s *GameState) Update() {
	var gold int
	fmt.Scan(&gold)
	s.gold = gold
	var enemyGold int
	fmt.Scan(&enemyGold)
	s.enemyGold = enemyGold
	// roundType: a positive value will show the number of heroes that await a command
	var roundType int
	fmt.Scan(&roundType)
	s.roundType = roundType
	var entityCount int
	fmt.Scan(&entityCount)
	s.units = make(Units, entityCount)
	for i := 0; i < entityCount; i++ {
		var unitId, team int
		var unitType string
		var x, y, attackRange, health, maxHealth, shield, attackDamage, movementSpeed, stunDuration, goldValue, countDown1, countDown2, countDown3, mana, maxMana, manaRegeneration int
		var heroType string
		var isVisible bool
		var itemsOwned int
		fmt.Scan(&unitId, &team, &unitType, &x, &y, &attackRange, &health, &maxHealth, &shield, &attackDamage, &movementSpeed, &stunDuration, &goldValue, &countDown1, &countDown2, &countDown3, &mana, &maxMana, &manaRegeneration, &heroType, &isVisible, &itemsOwned)
		s.units[i] = &Unit{unitId, team, unitType, Point{x, y}, attackRange, health, maxHealth, shield, attackDamage, movementSpeed, stunDuration, goldValue, countDown1, countDown2, countDown3, mana, maxMana, manaRegeneration, heroType, isVisible, itemsOwned}
	}
}

func (s *GameState) Turn() string {
	if s.roundType < 0 {
		return Wait()
	}
	myHeroes := s.units.Heroes(s.team)
	enemyHeroes := s.units.Heroes(s.team + 1%2)
	actions := make([]string, len(myHeroes))
	for _, hero := range myHeroes {
		target := enemyHeroes.First()
		action := AttackNearest("HERO")
		distance := hero.pos.Distance(target.pos)
		myArrival := arrivalTime(distance, hero.movementSpeed)
		enemyArrival := arrivalTime(distance, target.movementSpeed)
		debug("distance: %d, me: %v, them: %v\n", distance, myArrival, enemyArrival)
		if distance > hero.attackRange {
			action = MoveAttack(target.pos, target)
		}
		actions = append(actions, action)
	}
	return strings.Join(actions, "\n")
}

func main() {
	state := NewGame()
	for {
		state.Update()
		fmt.Println(state.Turn())
	}
}
