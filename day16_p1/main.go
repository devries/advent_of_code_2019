package main

import (
	"fmt"
	"strconv"
)

const input = "59754835304279095723667830764559994207668723615273907123832849523285892960990393495763064170399328763959561728553125232713663009161639789035331160605704223863754174835946381029543455581717775283582638013183215312822018348826709095340993876483418084566769957325454646682224309983510781204738662326823284208246064957584474684120465225052336374823382738788573365821572559301715471129142028462682986045997614184200503304763967364026464055684787169501819241361777789595715281841253470186857857671012867285957360755646446993278909888646724963166642032217322712337954157163771552371824741783496515778370667935574438315692768492954716331430001072240959235708"

func main() {
	digits, err := parseNumber(input)
	if err != nil {
		panic(fmt.Errorf("Error parsing input: %s", err))
	}

	o := RepeatPhase(digits, 100)
	fmt.Println(o[:8])

}

var basePattern = []int{0, 1, 0, -1}

func generatePattern(length int, repeat int) []int {
	r := make([]int, length+1)

	pt := 0
	for i := 0; i <= length; i++ {
		r[i] = basePattern[pt]
		if i%repeat == repeat-1 {
			pt = (pt + 1) % len(basePattern)
		}
	}

	return r[1:]
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

func FFTPhase(indigits []int) []int {
	l := len(indigits)
	out := make([]int, l)

	for i := 0; i < l; i++ {
		pat := generatePattern(l, i+1)

		digit := 0
		for j := 0; j < l; j++ {
			digit += indigits[j] * pat[j]
		}
		digit = Abs(digit) % 10
		out[i] = digit
	}

	return out
}

func RepeatPhase(indigits []int, repeat int) []int {
	out := make([]int, len(indigits))
	copy(out, indigits)

	for i := 0; i < repeat; i++ {
		out = FFTPhase(out)
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
