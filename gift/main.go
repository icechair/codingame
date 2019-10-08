package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
func debugln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

type Budget []int

func (b Budget) String() string {
	out := make([]string, len(b))
	for i, e := range b {
		out[i] = strconv.Itoa(e)
	}
	return strings.Join(out, "\n")
}

func (b Budget) Sum() int {
	sum := 0
	for _, e := range b {
		sum += e
	}
	return sum
}

func (b Budget) Mean() int {
	return b.Sum() / len(b)
}

func (b Budget) Median() int {
	m := len(b) / 2
	if len(b)%2 == 0 {
		return (b[m-1] + b[m]) / 2
	}
	return b[m]
}

func main() {
	var N int
	fmt.Scan(&N)

	var C int
	fmt.Scan(&C)
	budget := make(Budget, N)
	for i := 0; i < N; i++ {
		var B int
		fmt.Scan(&B)
		budget[i] = B
	}
	sort.Ints(budget)
	if budget.Sum() < C {
		fmt.Println("IMPOSSIBLE") // Write answer to stdout
		return
	}

	parts := make(Budget, N)
	for i, b := range budget {
		mean := C / (len(budget) - i)
		debugln(len(budget)-i, C, mean)
		if b < mean {
			mean = b
		}
		C -= mean
		parts[i] = mean
	}

	// debugln(budget)
	fmt.Println(parts)
}
