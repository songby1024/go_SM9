package main

import (
	"fmt"
	"math"
	"strconv"
)

func To10(num int) int {
	var s int
	str := strconv.Itoa(num)
	mi := len(str) - 1
	for _, x := range str {
		val, _ := strconv.Atoi(string(x))
		s += val * int(math.Pow(float64(16), float64(mi)))
		mi--
	}
	return s
}

func main() {
	start := 10
	for {
		if To10(start)%start == 0 {
			fmt.Println("res", start, To10(start))
			break
		}
		start++
	}

}
