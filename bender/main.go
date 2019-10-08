package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//import "strings"
//import "strconv"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
func debug(format string, a ...interface{}) {
	// 	fmt.Fprintf(os.Stderr, format, a...)
}

const (
	EMPTY    = rune(' ')
	START    = rune('@')
	END      = rune('$')
	WALL     = rune('#')
	BOX      = rune('X')
	SOUTH    = rune('S')
	EAST     = rune('E')
	NORTH    = rune('N')
	WEST     = rune('W')
	INVERTER = rune('I')
	BEER     = rune('B')
	TELEPORT = rune('T')
)

type Area [][]rune

func (a Area) Value(p Point) rune {
	return a[p.Row][p.Col]
}
func (a *Area) Update(p Point, r rune) {
	(*a)[p.Row][p.Col] = r
}

type Point struct {
	Row int
	Col int
}

func (p Point) Add(direction rune) Point {
	switch direction {
	case SOUTH:
		p.Row++
	case EAST:
		p.Col++
	case NORTH:
		p.Row--
	case WEST:
		p.Col--
	}
	return p
}

type Bender struct {
	Point
	Direction rune
	Breaker   bool
	Inverted  bool
}

func (b *Bender) IsPassable(field rune) bool {
	if field == WALL {
		return false
	}
	if field == BOX && !b.Breaker {
		return false
	}
	return true
}

func (b *Bender) NextValidDirection(area Area) rune {
	d := b.Direction
	p := b.Point.Add(d)
	r := area.Value(p)
	if b.IsPassable(r) {
		return b.Direction
	}
	values := [][2]rune{
		[2]rune{SOUTH, area.Value(b.Point.Add(SOUTH))},
		[2]rune{EAST, area.Value(b.Point.Add(EAST))},
		[2]rune{NORTH, area.Value(b.Point.Add(NORTH))},
		[2]rune{WEST, area.Value(b.Point.Add(WEST))},
	}
	if b.Inverted {
		values = [][2]rune{
			[2]rune{WEST, area.Value(b.Point.Add(WEST))},
			[2]rune{NORTH, area.Value(b.Point.Add(NORTH))},
			[2]rune{EAST, area.Value(b.Point.Add(EAST))},
			[2]rune{SOUTH, area.Value(b.Point.Add(SOUTH))},
		}
	}

	for _, kv := range values {
		if b.IsPassable(kv[1]) {
			return kv[0]
		}
	}
	debug("this shoudlnt happen %#v\n", b)
	return b.Direction
}

func (b *Bender) Turn(area *Area, teleport map[Point]Point) Bender {
	r := area.Value(b.Point)
	debug("%s %v %s\n", string(b.Direction), b.Point, string(r))
	switch r {
	case SOUTH, WEST, EAST, NORTH:
		b.Direction = r
	case BEER:
		b.Breaker = !b.Breaker
	case INVERTER:
		b.Inverted = !b.Inverted
	case TELEPORT:
		b.Point = teleport[b.Point]
	case BOX:
		area.Update(b.Point, EMPTY)
	}
	b.Direction = b.NextValidDirection(*area)

	p := b.Point.Add(b.Direction)
	out := *b
	b.Point = p
	return out
}

type History []Bender

func (h History) String() string {
	out := make([]string, len(h))
	for idx, b := range h {
		switch b.Direction {
		case SOUTH:
			out[idx] = "SOUTH"
		case WEST:
			out[idx] = "WEST"
		case NORTH:
			out[idx] = "NORTH"
		case EAST:
			out[idx] = "EAST"
		}
	}
	return strings.Join(out, "\n")
}

func (h History) Exists(a Bender) bool {
	for _, b := range h {
		if a == b {
			return true
		}
	}
	return false
}

func (h History) Print() string {
	var out []string
	for _, e := range h {
		out = append(out, fmt.Sprintf("%#v", e))
	}
	return strings.Join(out, "\n")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var L, C int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &L, &C)
	area := make(Area, L)
	b := &Bender{Point{0, 0}, SOUTH, false, false}
	tps := make([]Point, 0)
	var exit Point
	for row := 0; row < L; row++ {
		scanner.Scan()
		line := scanner.Text()
		for col, c := range line {
			area[row] = append(area[row], c)
			switch {
			case c == START:
				b.Row = row
				b.Col = col
			case c == TELEPORT:
				tps = append(tps, Point{row, col})
			case c == END:
				exit.Col = col
				exit.Row = row
			}
		}
		debug("%v\n", line)
	}
	teleport := make(map[Point]Point)
	if len(tps) == 2 {
		teleport[tps[0]] = tps[1]
		teleport[tps[1]] = tps[0]
	}
	debug("bender: %#v\n", b)
	debug("teleport: %#v\n", len(teleport))
	debug("exit: %#v\n", exit)
	history := make(History, 0)
	n := 0
	for n < 1000 && b.Point != exit {
		history = append(history, b.Turn(&area, teleport))
		n++
	}
	debug("history: %s\n", history.Print())
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	if b.Point != exit {
		fmt.Println("LOOP")
	} else {
		fmt.Println(history)
	}
}
