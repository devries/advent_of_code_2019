package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const ncards = 10007
const card = 2019

func main() {
	f, err := os.Open("input.txt")
	// f, err := os.Open("simple.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open input file"))
	}
	defer f.Close()

	dealFunctions, instructions, err := parseInput(f)
	if err != nil {
		panic(fmt.Errorf("Error reading file: %s", err))
	}

	c := card
	for i, df := range dealFunctions {
		c = df(c)
		fmt.Printf("%s : %d\n", instructions[i], c)
	}

	fmt.Printf("Location of card %d: %d\n", card, c)
}

type reindex func(int) int

func deal(n int) reindex {
	return func(i int) int {
		result := n - i - 1
		return result
	}
}

func dealCut(n int, cut int) reindex {
	if cut < 0 {
		cut = n + cut
	}

	return func(i int) int {
		result := i - cut
		if result < 0 {
			result = n + result
		}
		return result
	}
}

func dealIncr(n int, incr int) reindex {
	return func(i int) int {
		result := (incr * i) % n
		return result
	}
}

func parseInput(r io.Reader) ([]reindex, []string, error) {
	var result []reindex
	var instructions []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		instructions = append(instructions, line)
		if strings.HasPrefix(line, "deal into") {
			result = append(result, deal(ncards))
		} else if strings.HasPrefix(line, "cut") {
			var cut int
			_, err := fmt.Sscanf(line, "cut %d", &cut)
			if err != nil {
				return nil, nil, err
			}
			result = append(result, dealCut(ncards, cut))
		} else {
			var incr int
			_, err := fmt.Sscanf(line, "deal with increment %d", &incr)
			if err != nil {
				return nil, nil, err
			}
			result = append(result, dealIncr(ncards, incr))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return result, instructions, nil
}
