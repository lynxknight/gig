package distance

import (
	"testing"
)

func TestMin(t *testing.T) {
	m := min(3, 4, 5, 6, 8, -9)
	if m != -9 {
		t.Errorf("Mintest wrong =(: %v", m)
	}
}

// Levenshtein distance test
type ldt struct {
	s1, s2    string
	edistance int
}

var testCasesLDT = []ldt{
	ldt{s1: "A", s2: "A", edistance: 0},
	ldt{s1: "A", s2: "B", edistance: 1},
	ldt{s1: "A", s2: "AB", edistance: 1},
	ldt{s1: "AA", s2: "AB", edistance: 1},
	ldt{s1: "AAA", s2: "AAB", edistance: 1},
	ldt{s1: "AAA", s2: "B", edistance: 3},
	ldt{s1: "AAAA", s2: "B", edistance: 4},
	ldt{s1: "kitten", s2: "sitting", edistance: 3},
	ldt{s1: "saturday", s2: "sunday", edistance: 3},
	ldt{s1: "exponential", s2: "polynomial", edistance: 6},
}

// calc Levenshtein score test
type clst struct {
	s1, s2 string
	score  Score
}

func TestDistanceIterative(t *testing.T) {
	for _, test := range testCasesLDT {
		if d := calcLeventshteinDistance(test.s1, test.s2, levenThreshold); d != test.edistance {
			t.Errorf(
				"'%v' and '%v' should have distance of %v, got: %v",
				test.s1, test.s2, test.edistance, d,
			)
		}
		t.Logf(
			"[PASSED] s1=%v s2=%v edistance=%v",
			test.s1, test.s2, test.edistance,
		)
	}
}

var testCasesCLST = []clst{
	clst{s1: "feature/data-v3-reload-nothing", s2: "data_v3", score: Score{1, 8, 15}},
}

func TestCalcLevenshteinScore(t *testing.T) {
	for _, test := range testCasesCLST {
		if s := calcLevenshteinScore(test.s1, test.s2); s != test.score {
			t.Errorf(
				"'%v' and '%v' should have score of %v, got: %v",
				test.s1, test.s2, test.score, s,
			)
		} else {
			t.Logf(
				"[PASSED] s1=%v s2=%v score=%v",
				test.s1, test.s2, test.score,
			)
		}
	}
}

func TestGetScore(t *testing.T) {
	destScore := Score{1, 0, 7}
	if s := GetScore("feature/1/2/3", "featyre"); s != destScore {
		t.Errorf("%v", s)
	}
}

func TestGetScoreStrShorterThanTarget(t *testing.T) {
	destScore := Score{1, 0, 2}
	if s := GetScore("12", "123"); s != destScore {
		t.Errorf("%v", s)
	}
}

func TestGetScoreStrShorterThanTargetButOverThreshold(t *testing.T) {
	destScore := Score{2, 0, 2}
	if s := GetScore("12", "1234"); s != destScore {
		t.Errorf("\nWant: %v\nHave: %v", destScore, s)
	}
}

func TestGetScoreStrShorterThanTargetButMuchOverThreshold(t *testing.T) {
	destScore := Score{6, 0, 2}
	if s := GetScore("12", "12345678"); s != destScore {
		t.Errorf("\nWant: %v\nHave: %v", destScore, s)
	}
}

func TestGetScoreIsCaseInsensitive(t *testing.T) {
	s1 := GetScore("abcd", "aBcDe")
	s2 := GetScore("abcd", "abcde")
	if s1 != s2 {
		t.Errorf("\nHave (s1: %v; s2: %v; s1 != s2)\nWant: s1 == s2", s1, s2)
	}
}

// In case of troubles with some test
// func TestSpecificCase(t *testing.T) {
// 	test := levTest{s1: "saturday", s2: "sunday", edistance: 3}
// 	if d := Distance(test.s1, test.s2); d != test.edistance {
// 		t.Errorf(
// 			"'%v' and '%v' should have distance of %v, got: %v",
// 			test.s1, test.s2, test.edistance, d,
// 		)
// 	}
// }
