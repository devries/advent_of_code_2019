package main

import (
	"fmt"
	"strconv"
)

func main() {
	input := "171309-643603"
	var start int
	var stop int

	_, err := fmt.Sscanf(input, "%d-%d", &start, &stop)
	if err != nil {
		panic(fmt.Errorf("Unable to scan input"))
	}

	counter := 0
	for i := start; i <= stop; i++ {
		if ValidateCode(i) {
			counter++
		}
	}

	fmt.Printf("There are %d valid codes\n", counter)
}

func ValidateCode(code int) bool {
	s := strconv.Itoa(code)

	characters := []rune(s)
	if len(characters) != 6 {
		return false
	}

	seen := make(map[rune]bool)
	double_detected := false
	startVal := rune(0)
	for _, v := range characters {
		if v < startVal {
			return false
		}
		startVal = v
		if seen[v] {
			double_detected = true
		}
		seen[v] = true
	}

	return double_detected
}
