package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open input file"))
	}
	defer f.Close()

	c, err := parseInput(f)
	if err != nil {
		panic(fmt.Errorf("Error parsing input"))
	}

	totalOre := 1000000000000
	shopping := NewShoppingList()
	fuel := 0
	for {
		need(1, "FUEL", c, shopping)
		if shopping["ORE"].BatchesMade > totalOre {
			break
		}
		fuel++
		if fuel%10000 == 0 {
			fmt.Println(fuel, shopping["ORE"].BatchesMade)
		}
	}

	fmt.Printf("%d fuel produced\n", fuel)
}

type Recipe struct {
	Batchsize   int
	Ingredients map[string]int
}

type ShoppingItem struct {
	Required    int
	BatchesMade int
}

type Cookbook map[string]Recipe

type ShoppingList map[string]ShoppingItem

func NewShoppingList() ShoppingList {
	r := make(map[string]ShoppingItem)
	return r
}

func NewCookbook() Cookbook {
	r := make(map[string]Recipe)

	r["ORE"] = Recipe{1, nil}

	return r
}

func (c Cookbook) Add(name string, r Recipe) {
	c[name] = r
}

func parseLine(l string, c Cookbook) error {
	middle := strings.Index(l, "=>")

	input := l[0:middle]
	output := l[middle+2:]

	var batchsize int
	var product string
	n, err := fmt.Sscanf(output, "%d %s", &batchsize, &product)
	if n != 2 {
		return fmt.Errorf("Error scanning output: %s", output)
	} else if err != nil {
		return err
	}

	r := Recipe{batchsize, make(map[string]int)}
	for _, v := range strings.Split(input, ",") {
		var quantity int
		var ingredient string
		n, err := fmt.Sscanf(v, "%d %s", &quantity, &ingredient)
		if n != 2 {
			return fmt.Errorf("Error scanning ingredient: %s", v)
		} else if err != nil {
			return err
		}
		r.Ingredients[ingredient] = quantity
	}

	c.Add(product, r)

	return nil
}

func parseInput(r io.Reader) (Cookbook, error) {
	c := NewCookbook()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		err := parseLine(line, c)
		if err != nil {
			return nil, fmt.Errorf("Error scanning line %s: %s", line, err)
		}
	}

	return c, nil
}

// My mom told me not to recurse at the dinner table, so I did this in my office.
func need(quantity int, ingredient string, c Cookbook, l ShoppingList) {
	_, ok := l[ingredient]
	if !ok {
		l[ingredient] = ShoppingItem{quantity, 0}
	} else {
		r := l[ingredient]
		r.Required += quantity
		l[ingredient] = r
	}

	ci := c[ingredient]
	if l[ingredient].Required > l[ingredient].BatchesMade*ci.Batchsize {
		nBatchesRequired := l[ingredient].Required / ci.Batchsize
		if l[ingredient].Required%ci.Batchsize != 0 {
			nBatchesRequired++
		}
		nBatchesToMake := nBatchesRequired - l[ingredient].BatchesMade
		for k, v := range c[ingredient].Ingredients {
			need(v*nBatchesToMake, k, c, l)
		}
		b := l[ingredient]
		b.BatchesMade += nBatchesToMake
		l[ingredient] = b
	}
}
