package main

import (
	"fmt"
	"strconv"
	"time"
)

const input = "59754835304279095723667830764559994207668723615273907123832849523285892960990393495763064170399328763959561728553125232713663009161639789035331160605704223863754174835946381029543455581717775283582638013183215312822018348826709095340993876483418084566769957325454646682224309983510781204738662326823284208246064957584474684120465225052336374823382738788573365821572559301715471129142028462682986045997614184200503304763967364026464055684787169501819241361777789595715281841253470186857857671012867285957360755646446993278909888646724963166642032217322712337954157163771552371824741783496515778370667935574438315692768492954716331430001072240959235708"

func main() {
	digits, err := parseNumber(input)
	// digits, err := parseNumber("03036732577212944063491565474664")
	// digits, err := parseNumber("02935109699940807407585447034323")
	// digits, err := parseNumber("03081770884921959731165446850517")
	if err != nil {
		panic(fmt.Errorf("Error parsing input: %s", err))
	}

	out := Expansion(digits)

	offset := 0
	for i := 0; i < 7; i++ {
		offset *= 10
		offset += out[i]
	}
	fmt.Println(offset)
	out = RepeatFFT(out, 100, offset, offset+8)
	fmt.Println(out[offset : offset+8])
}

func parseNumber(s string) ([]int, error) {
	numbers := []rune(s)
	r := make([]int, len(numbers))

	for i := 0; i < len(numbers); i++ {
		v, err := strconv.Atoi(string(numbers[i]))
		if err != nil {
			return nil, err
		}
		r[i] = v
	}

	return r, nil
}

func Expansion(indigits []int) []int {
	out := make([]int, 10000*len(indigits))

	for i := 0; i < 10000; i++ {
		copy(out[i*len(indigits):], indigits)
	}

	return out
}

func RepeatFFT(indigits []int, repeat int, start, stop int) []int {
	l := len(indigits)
	out := make([]int, l)
	temp := make([]int, l)
	copy(temp[start:], indigits[start:])

	tstart := time.Now()
	var digit int

	for i := 0; i < repeat; i++ {
		elapsed := time.Since(tstart)
		fmt.Printf("r(%d) -- %s elapsed.\n", i, elapsed)
		sum := 0
		for k := l - 1; k >= Min(l/2, start); k-- {
			sum += temp[k]
			digit = Abs(sum) % 10
			out[k] = digit
		}
		for k := l/2 - 1; k >= Min(l/4, start); k-- {
			sum += temp[k]
			sum -= temp[2*k]
			sum -= temp[2*k+1]
			digit = Abs(sum) % 10
			out[k] = digit
		}
		for k := l/4 - 1; k >= start; k++ {
			sum += temp[k]
			sum -= temp[2*k]
			sum -= temp[2*k+1]
			sum -= temp[3*k]
			sum -= temp[3*k+1]
			sum -= temp[3*k+2]
			digit = Abs(sum) % 10
			out[k] = digit
		}
		copy(temp[start:], out[start:])
	}

	return out
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}
