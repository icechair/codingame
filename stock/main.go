package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
func debugln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func Max(a ...int) int {
	max := 0
	for _, x := range a {
		if x > max {
			max = x
		}
	}
	return max
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var n int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)

	scanner.Scan()
	inputs := strings.Split(scanner.Text(), " ")
	stock := make([]int, n)
	drop := 0
	max := 0
	for i := 0; i < n; i++ {
		v, _ := strconv.Atoi(inputs[i])
		stock[i] = v
		max = Max(max, v)
		drop = Max(drop, max-v)
	}
	debugln(stock)
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println(drop * -1) // Write answer to stdout
}
