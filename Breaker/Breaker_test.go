package Breaker

import "testing"

func TestNewBreaker(t *testing.T) {
	got := NewBreakerByte([]byte{1, 2, 3})
	if len(got.CipherStream) != 3 {
		t.Error("Didn't get back a struct.")
	}
}

//func TestScore(t *testing.T){
//	testScoreFunc(score, t)
//	// All e's should score low
//
//}

func TestNewScore(t *testing.T) {
	testScoreFunc(score, t)
}

func testScoreFunc(f func([]byte) float64, t *testing.T) {
	lowThreshold := 35.0
	expectLow := []string{
		"eeeeeeeeeeeeeeeeeeee",
		"eEeEeEeEeEeEeEeEeEeE",
		"EEEEEEEEEEEEEEEEEEEE",
		// All spaces/punctuation should score low
		"                    ",
		".,;[#%^!$*&&$!$%><?:",
	}
	// Actual English should score high
	highThreshold := 50.0
	expectHigh := []string{
		"This is an english sentence",
		"Verily I was perplex'd by the code!",
		// Letters that are scrambled english should also score high
		"ylirVeIswa'plefjicdea       ",
	}
	for _, v := range expectLow {
		s := f([]byte(v))
		if s < 0 || s > maxScore {
			t.Errorf("Out of Bounds Error - Over 100 or Under 0 (%s) scored %f", v, s)
		}
		if s > lowThreshold {
			t.Errorf("Scoring error over threshold %f (%s) scored %f ", lowThreshold, v, s)
		}
	}
	for _, v := range expectHigh {
		s := f([]byte(v))
		if s < 0 || s > maxScore {
			t.Errorf("Out of Bounds Error - Over 100 or Under 0 (%s) scored %f", v, s)
		}
		if s < highThreshold {
			t.Errorf("Scoring error under threshold %f (%s) scored %f ", highThreshold, v, s)
		}
	}
}
