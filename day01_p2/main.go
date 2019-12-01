package main

import (
	"bufio"
	"fmt"
	"os"
)

var input_filename string = "input.txt"

func main() {
	f, err := os.Open(input_filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	fuelSum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var mass int
		_, err := fmt.Sscanf(scanner.Text(), "%d", &mass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		for {
			fuel := moduleFuel(mass)
			if fuel <= 0 {
				break
			}
			fuelSum += fuel
			mass = fuel
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total Fuel: %d\n", fuelSum)
}

func moduleFuel(m int) int {
	return m/3 - 2
}
