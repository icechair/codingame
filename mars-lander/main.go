package main

import (
	"fmt"
	"math"
	"os"
)

//Debug prints to stderr for codingame
func Debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", format), a...)
}

const GRAVITY = 3.711

type Point struct {
	x float64
	y float64
}
type Surface []*Point

func (s Surface) GetLanding() Landing {
	var landing = Landing{}
	var prev *Point
	for idx, p := range s {
		Debug("iteration(%d): prev: %#v, p%#v", idx, prev, p)
		if idx == 0 {
			prev = p
			continue
		}
		if prev.y == p.y {
			landing.start = prev
			landing.end = p
			return landing
		} else {
			prev = p
		}
	}
	return landing
}

type Landing struct {
	start *Point
	end   *Point
}

type Vector struct {
	angle  float64
	length float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{
		angle:  math.Atan(y / x),
		length: math.Sqrt(x*x + y*y),
	}
}

type Rocket struct {
	pos    Point
	speed  Vector
	fuel   int
	rotate int
	power  int
}

func (r *Rocket) Update(x, y, hSpeed, vSpeed, rotate, power int) {
	r.pos.x = float64(x)
	r.pos.y = float64(y)
	r.speed = *NewVector(float64(hSpeed), float64(vSpeed))
	r.rotate = rotate
	r.power = power

}
func (r *Rocket) Turn(x, y, hSpeed, vSpeed, rotate, power int) string {
	r.Update(x, y, hSpeed, vSpeed, rotate, power)

	return "-20 3"
}

func FirstInput() Surface {
	var surfaceN int
	fmt.Scan(&surfaceN)
	surface := make(Surface, surfaceN)
	for i := 0; i < surfaceN; i++ {
		// landX: X coordinate of a surface point. (0 to 6999)
		// landY: Y coordinate of a surface point. By linking all the points together in a sequential fashion, you form the surface of Mars.
		var landX, landY int
		fmt.Scan(&landX, &landY)
		surface[i] = &Point{x: float64(landX), y: float64(landY)}
		Debug("%#v", surface[i])
	}
	return surface
}
func main() {
	// surfaceN: the number of points used to draw the surface of Mars.
	surface := FirstInput()
	Debug("%#v", surface)
	landing := surface.GetLanding()
	Debug("%v, %v", landing.start, landing.end)
	player := &Rocket{}
	for {
		// hSpeed: the horizontal speed (in m/s), can be negative.
		// vSpeed: the vertical speed (in m/s), can be negative.
		// fuel: the quantity of remaining fuel in liters.
		// rotate: the rotation angle in degrees (-90 to 90).
		// power: the thrust power (0 to 4).
		var X, Y, hSpeed, vSpeed, fuel, rotate, power int
		fmt.Scan(&X, &Y, &hSpeed, &vSpeed, &fuel, &rotate, &power)

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// rotate power. rotate is the desired rotation angle. power is the desired thrust power.
		fmt.Println(player.Turn(X, Y, hSpeed, vSpeed, rotate, power))
	}
}
