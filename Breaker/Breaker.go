package Breaker

import (
	"encoding/base64"
	"encoding/hex"
	"math"
	"sort"
)

const maxScore = 114.5

type FrequencyResults []FrequencyScore

func (f FrequencyResults) Len() int {
	return len(f)
}

func (f FrequencyResults) Less(i int, j int) bool {
	x := f[i].BestScore
	y := f[j].BestScore
	if x > y {
		return true
	}
	return false
}

func (f FrequencyResults) Swap(i int, j int) {
	f[i], f[j] = f[j], f[i]
}

func (results FrequencyResults) GetBestFrequencyScore() FrequencyScore {
	sort.Sort(results)
	return results[0]
}

type Breaker struct {
	Results      FrequencyResults
	CipherStream []byte
}

type FrequencyScore struct {
	BestResult   string
	BestKey      string
	BestScore    float64
	BestKeyBytes []byte
	Language     string
}

func (f FrequencyScore) SetLanguage(l string) {
	f.Language = l
}
func scoreOld(b []byte) float64 {
	score := 0.00
	for _, v := range b {
		if v < 32 {
			score -= 100
		}
		currentLetter := string(v)
		score += letterFrequency[currentLetter]
	}
	return score
}

func NewBreakerByte(cipherStream []byte) Breaker {
	return Breaker{CipherStream: cipherStream}
}

func NewBreakerHex(encodedHex string) Breaker {
	b, _ := hex.DecodeString(encodedHex)
	return Breaker{CipherStream: b}
}

func NewBreakerBase64(encoded string) Breaker {
	b, _ := base64.StdEncoding.DecodeString(encoded)
	return Breaker{CipherStream: b}
}

func NewBreaker(cipherText string) Breaker {
	return NewBreakerByte([]byte(cipherText))
}

func (b Breaker) BreakIt() FrequencyResults {
	return b.TryRange(32, 122)
}

func (b Breaker) TryKey(key byte) FrequencyScore {
	keyBytes := []byte{key}
	xorResult := Xor(b.CipherStream, keyBytes)
	resultScore := score(xorResult)
	return FrequencyScore{BestResult: string(xorResult), BestScore: resultScore, BestKey: string(key), BestKeyBytes: keyBytes} //TODO: IMPLEMENT STUB
}

func Xor(b1, b2 []byte) []byte {
	b2len := len(b2)
	result := []byte{}
	for i, v := range b1 {

		result = append(result, v^b2[i%b2len])

	}
	return result
}

func (b Breaker) TryRange(startAt byte, endAt byte) FrequencyResults {
	results := FrequencyResults{}
	for j := startAt; j < endAt; j++ {
		f := b.TryKey(j)
		results = append(results, f)
	}
	sort.Sort(results)
	return results
}

func (b Breaker) GetBestResult() FrequencyScore {

	return b.Results.GetBestFrequencyScore()
}

func score(b []byte) float64 {
	// First, count all the letters.
	letterCount := map[string]float64{}
	currentScore := maxScore
	totalCharacters := float64(len(b))
	for _, v := range b {
		if v < 32 {
			switch v {
			case 10:
				// Treat newlines as spaces.
				letterCount[" "]++
			default:
				letterCount["noise"]++
			}

		} else if v == 32 {
			letterCount[" "]++
		} else if v < 64 {
			letterCount["punctuation"]++
		} else if v < 91 {
			letterCount[string(v+32)]++
		} else if v < 97 {
			letterCount["punctuation"]++
		} else if v < 123 {
			letterCount[string(v)]++
		} else {
			letterCount["punctuation"]++
		}

	}
	// Take the difference between the observed frequency and expected frequency
	var noisePenalty float64
	for key, value := range letterCount {
		observedFrequency := value / totalCharacters * 100
		frequencyDifference := math.Abs(letterFrequency[key] - observedFrequency)
		switch key {
		case "noise":
			noisePenalty = frequencyDifference
		default:
			currentScore -= frequencyDifference
		}
	}
	currentScore /= math.Exp(noisePenalty)

	// Square the difference

	// Sum all of the squared differences

	// Subtract that from the max score, which is 100**2, or 10000

	// Calculate the expected difference of the noise

	// Divide the score by e**squared difference of the noise.
	return currentScore

}

var letterFrequency = map[string]float64{
	"E":           12.02,
	"e":           12.02,
	"T":           9.10,
	"t":           9.10,
	"A":           8.12,
	"a":           8.12,
	"O":           7.68,
	"o":           7.68,
	"I":           7.31,
	"i":           7.31,
	"N":           6.95,
	"n":           6.95,
	"S":           6.28,
	"s":           6.28,
	"R":           6.02,
	"r":           6.02,
	"H":           5.92,
	"h":           5.92,
	"D":           4.32,
	"d":           4.32,
	"L":           3.98,
	"l":           3.98,
	"U":           2.88,
	"u":           2.88,
	"C":           2.71,
	"c":           2.71,
	"M":           2.61,
	"m":           2.61,
	"F":           2.30,
	"f":           2.30,
	"Y":           2.11,
	"y":           2.11,
	"W":           2.09,
	"w":           2.09,
	"G":           2.03,
	"g":           2.03,
	"P":           1.82,
	"p":           1.82,
	"B":           1.49,
	"b":           1.49,
	"V":           1.11,
	"v":           1.11,
	"K":           0.69,
	"k":           0.69,
	"X":           0.17,
	"x":           0.17,
	"Q":           0.11,
	"q":           0.11,
	"J":           0.10,
	"j":           0.10,
	"Z":           0.07,
	"z":           0.07,
	" ":           14.3,
	"punctuation": 0.2,
	"noise":       0.0,
	".":           1.00,
	"'":           1.00,
	"&":           -10,
	"$":           -10,
	"#":           -10,
	"@":           -10,
	"*":           -10,
	"(":           -10,
	")":           -10,
	"[":           -10,
	"]":           -10,
	"{":           -10,
	"}":           -10,
	"`":           -10,
}
