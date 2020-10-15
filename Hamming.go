package main

import "encoding/base64"

type hamDistance struct {
	Input      string
	InputBytes []byte
}

type hamKeyScore struct {
	Keysize  int
	HamScore float64
	Distance int
}

func NewHammingBase64(s string) hamDistance {
	decoded, _ := base64.StdEncoding.DecodeString(s)
	return hamDistance{Input: string(decoded), InputBytes: decoded}
}

func (input hamDistance) GetBestKeySize(min int, max int) int {
	bestHamming := 1000
	bestHammingKeysize := 0
	for i := min; i < max; i++ {
		s1 := input.InputBytes[0:i]
		s2 := input.InputBytes[i:(2 * i)]
		currentHamming := hamming(string(s1), string(s2))
		currentHamming /= i
		if currentHamming < bestHamming {
			bestHamming = currentHamming
			bestHammingKeysize = i
		}
	}
	return bestHammingKeysize
}

func (input hamDistance) Pivot(x int) [][]byte {
	s := [][]byte{}
	for i := 0; i < x; i++ {
		s = append(s, []byte{})
	}
	for i, v := range input.InputBytes {
		s[i%x] = append(s[i%x], v)
	}
	return s
}
func hamming(s1, s2 string) int {
	s1bits := []uint8(s1)
	s2bits := []uint8(s2)
	xoroutput := xor(s1bits, s2bits)
	count := 0
	for _, i := range xoroutput {
		count += lookupBits(i)
	}

	return count
}

func lookupBits(x uint8) int {
	count := 0
	count += answerKey[x]
	return count
}

var answerKey = [256]int{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8}
