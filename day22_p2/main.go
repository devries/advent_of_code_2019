package main

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
)

const ncards = 119315717514047
const card = 2020
const nshuffles = 101741582076661

/*
const ncards = 10007
const card = 2019
*/

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open input file"))
	}
	defer f.Close()

	lf, _, err := parseInput(f, ncards)
	if err != nil {
		panic(fmt.Errorf("Error reading file: %s", err))
	}

	lf = lf.Repeat(big.NewInt(nshuffles))
	r := lf.Inverse(big.NewInt(card))

	fmt.Printf("Card at %s is in position %d\n", r, card)
}

type LinearFunction struct {
	A   *big.Int
	B   *big.Int
	Mod *big.Int
}

func NewLinearFunction(a, b, mod int64) LinearFunction {
	return LinearFunction{big.NewInt(a), big.NewInt(b), big.NewInt(mod)}
}

func (lf LinearFunction) Calc(x *big.Int) *big.Int {
	r := &big.Int{}
	r.Mod(r.Add(r.Mul(lf.A, x), lf.B), lf.Mod)

	if r.Sign() < 0 {
		r.Add(r, lf.Mod)
	}

	return r
}

func (lf LinearFunction) Of(lfo LinearFunction) LinearFunction {
	newA := &big.Int{}
	newB := &big.Int{}
	newMod := &big.Int{}

	newA.Mod(newA.Mul(lf.A, lfo.A), lf.Mod)
	newB.Mod(newB.Add(newB.Mul(lf.A, lfo.B), lf.B), lf.Mod)
	newMod = lf.Mod

	return LinearFunction{newA, newB, newMod}
}

func (lf LinearFunction) Inverse(x *big.Int) *big.Int {
	r := &big.Int{}
	d := &big.Int{}

	r.Sub(x, lf.B)

	d = ModularDivision(r, lf.A, lf.Mod)
	return d
}

func (lf LinearFunction) Repeat(n *big.Int) LinearFunction {
	// Answer for N repeats is
	// f = a^{N}x + (a^{N-1}+a^{N-2}+ ... + a^{2}+a+1)*b
	// Geometric series yields
	// f = a^{N}x + ((a^{N}-1)/(a-1))*b
	newA := &big.Int{}
	newB := &big.Int{}
	newMod := &big.Int{}

	newA.Exp(lf.A, n, lf.Mod)

	numerator := &big.Int{}
	denominator := &big.Int{}

	numerator.Mul(numerator.Sub(newA, big.NewInt(1)), lf.B)
	denominator.Sub(lf.A, big.NewInt(1))

	newB = ModularDivision(numerator, denominator, lf.Mod)
	newMod = lf.Mod
	return LinearFunction{newA, newB, newMod}
}

func ModularDivision(numerator, denominator, mod *big.Int) *big.Int {
	modInverse := &big.Int{}
	result := &big.Int{}

	modInverse = modInverse.ModInverse(denominator, mod)
	result.Mod(result.Mul(numerator.Mod(numerator, mod), modInverse.Mod(modInverse, mod)), mod)

	return result
}

func identityFunction(mod int64) LinearFunction {
	return NewLinearFunction(1, 0, mod)
}

func cutFunction(cut, mod int64) LinearFunction {
	return NewLinearFunction(1, -cut, mod)
}

func dealFunction(mod int64) LinearFunction {
	return NewLinearFunction(-1, mod-1, mod)
}

func dealSkipFunction(skip, mod int64) LinearFunction {
	return NewLinearFunction(skip, 0, mod)
}

func parseInput(r io.Reader, ncards int64) (LinearFunction, []string, error) {
	result := identityFunction(ncards)
	var instructions []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		instructions = append(instructions, line)
		if strings.HasPrefix(line, "deal into") {
			result = dealFunction(ncards).Of(result)
		} else if strings.HasPrefix(line, "cut") {
			var cut int64
			_, err := fmt.Sscanf(line, "cut %d", &cut)
			if err != nil {
				return LinearFunction{}, nil, err
			}
			result = cutFunction(cut, ncards).Of(result)
		} else {
			var incr int64
			_, err := fmt.Sscanf(line, "deal with increment %d", &incr)
			if err != nil {
				return LinearFunction{}, nil, err
			}
			result = dealSkipFunction(incr, ncards).Of(result)
		}
	}
	if err := scanner.Err(); err != nil {
		return LinearFunction{}, nil, err
	}

	return result, instructions, nil
}
