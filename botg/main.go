package main

import (
	"fmt"
	"os"
	"strings"
)

func debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

const (
	NOTHING = iota
	BUSH
	SPAWN
	UNIT
	HERO
	TOWER
	GROOT
	VALKYRIE
	DEADPOOL
	IRONMAN
	DOCTOR_STRANGE
	HULK
)

const TURN_TIME = 1.0
const HERO_ATTACK_TIME = 0.1
const UNIT_ATTACK_TIME = 0.2

func attackTime(category, distance, attackRange int) float64 {
	time := UNIT_ATTACK_TIME
	if category == HERO {
		time = HERO_ATTACK_TIME
	}
	return time * float64(distance) / float64(attackRange)
}

func travelTime(distance, movementSpeed int) float64 {
	return float64(distance) / float64(movementSpeed)
}

func GetCategory(cat string) int {
	switch cat {
	default:
		return NOTHING
	case "BUSH":
		return BUSH
	case "SPAWN":
		return SPAWN
	case "UNIT":
		return UNIT
	case "HERO":
		return HERO
	case "TOWER":
		return TOWER
	case "GROOT":
		return GROOT
	case "VALKYRIE":
		return VALKYRIE
	case "DEADPOOL":
		return DEADPOOL
	case "IRONMAN":
		return IRONMAN
	case "DOCTOR_STRANGE":
		return DOCTOR_STRANGE
	case "HULK":
		return HULK
	}
}

func GetType(cat int) string {
	switch cat {
	default:
		return "WAIT"
	case VALKYRIE:
		return "VALKYRIE"
	case DEADPOOL:
		return "DEADPOOL"
	case IRONMAN:
		return "IRONMAN"
	case DOCTOR_STRANGE:
		return "DOCTOR_STRANGE"
	case HULK:
		return "HULK"
	}
}

type Point struct {
	x int
	y int
}

type Entity struct {
	category int
	pos      Point
	radius   int
}
type Entities []*Entity

func NewEntity(category string, x, y, radius int) *Entity {
	return &Entity{GetCategory(category), Point{x, y}, radius}
}

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
	category         int
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
	heroType         int
	isVisible        bool
	itemsOwned       int
}
type Units []*Unit

type Player struct {
	entities  Entities
	items     Items
	units     Units
	team      int
	gold      int
	enemyGold int
}

func (p *Player) Turn(roundType int) []string {
	actions := make([]string, roundType)
	for i, _ := range actions {
		actions[i] = "WAIT"
	}

	return actions
}

func main() {
	var myTeam int
	fmt.Scan(&myTeam)

	// bushAndSpawnPointCount: usefrul from wood1, represents the number of bushes and the number of places where neutral units can spawn
	var bushAndSpawnPointCount int
	fmt.Scan(&bushAndSpawnPointCount)
	entities := make(Entities, bushAndSpawnPointCount)
	for i := 0; i < bushAndSpawnPointCount; i++ {
		// entityType: BUSH, from wood1 it can also be SPAWN
		var entityType string
		var x, y, radius int
		fmt.Scan(&entityType, &x, &y, &radius)
		entities = append(entities, NewEntity(entityType, x, y, radius))
	}
	// itemCount: useful from wood2
	var itemCount int
	fmt.Scan(&itemCount)
	items := make(Items, itemCount)
	for i := 0; i < itemCount; i++ {
		// itemName: contains keywords such as BRONZE, SILVER and BLADE, BOOTS connected by "_" to help you sort easier
		// itemCost: BRONZE items have lowest cost, the most expensive items are LEGENDARY
		// damage: keyword BLADE is present if the most important item stat is damage
		// moveSpeed: keyword BOOTS is present if the most important item stat is moveSpeed
		// isPotion: 0 if it's not instantly consumed
		var itemName string
		var itemCost, damage, health, maxHealth, mana, maxMana, moveSpeed, manaRegeneration int
		var isPotion bool

		fmt.Scan(&itemName, &itemCost, &damage, &health, &maxHealth, &mana, &maxMana, &moveSpeed, &manaRegeneration, &isPotion)
		items = append(items, &Item{itemName, itemCost, damage, health, maxHealth, mana, maxMana, moveSpeed, manaRegeneration, isPotion})
	}

	info := &Player{entities, items, nil, myTeam, 0, 0}
	for {
		var gold int
		fmt.Scan(&gold)

		var enemyGold int
		fmt.Scan(&enemyGold)

		// roundType: a positive value will show the number of heroes that await a command
		var roundType int
		fmt.Scan(&roundType)

		var entityCount int
		fmt.Scan(&entityCount)
		units := make(Units, entityCount)
		for i := 0; i < entityCount; i++ {
			// unitType: UNIT, HERO, TOWER, can also be GROOT from wood1
			// shield: useful in bronze
			// stunDuration: useful in bronze
			// countDown1: all countDown and mana variables are useful starting in bronze
			// heroType: DEADPOOL, VALKYRIE, DOCTOR_STRANGE, HULK, IRONMAN
			// isVisible: 0 if it isn't
			// itemsOwned: useful from wood1
			var unitId, team int
			var unitType string
			var x, y, attackRange, health, maxHealth, shield, attackDamage, movementSpeed, stunDuration, goldValue, countDown1, countDown2, countDown3, mana, maxMana, manaRegeneration, itemsOwned int
			var heroType string
			var isVisible bool
			fmt.Scan(&unitId, &team, &unitType, &x, &y, &attackRange, &health, &maxHealth, &shield, &attackDamage, &movementSpeed, &stunDuration, &goldValue, &countDown1, &countDown2, &countDown3, &mana, &maxMana, &manaRegeneration, &heroType, &isVisible, &itemsOwned)
			units = append(units, &Unit{unitId, team, GetCategory(unitType), Point{x, y}, attackRange, health, maxHealth, shield, attackDamage, movementSpeed, stunDuration, goldValue, countDown1, countDown2, countDown3, mana, maxMana, manaRegeneration, GetCategory(heroType), isVisible, itemsOwned})
		}
		info.units = units
		info.gold = gold
		info.enemyGold = enemyGold
		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// If roundType has a negative value then you need to output a Hero name, such as "DEADPOOL" or "VALKYRIE".
		// Else you need to output roundType number of any valid action, such as "WAIT" or "ATTACK unitId"
		out := GetType(IRONMAN)
		if roundType > 0 {
			out = strings.Join(info.Turn(roundType), "\n")
		}
		fmt.Println(out)
	}
}
