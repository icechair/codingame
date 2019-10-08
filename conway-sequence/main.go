package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//import "os"

func Debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", format), a...)
}

func concat(a, b string) string {
	if len(a) > 0 {
		a += " "
	}
	return a + b
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
func describe(line string) string {
	digits := strings.Split(line, " ")
	n := 0
	tmp := ""
	out := ""
	for idx, digit := range digits {
		if idx == 0 {
			tmp = digit
		}
		if tmp == digit {
			n++
		} else {
			out = concat(out, fmt.Sprintf("%d %s", n, tmp))
			tmp = digit
			n = 1
		}
	}
	out = concat(out, fmt.Sprintf("%d %s", n, tmp))
	Debug("%s", out)
	return out
}
func main() {
	var R int
	fmt.Scan(&R)

	var L int
	fmt.Scan(&L)
	if L == 1 {
		fmt.Printf("%d\n", R)
		return
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	answer := describe(strconv.Itoa(R))
	for i := 1; i < L-1; i++ {
		answer = describe(answer)
	}
	fmt.Println(answer) // Write answer to stdout
}
