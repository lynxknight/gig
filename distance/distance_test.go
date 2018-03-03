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

type levTest struct {
	s1, s2    string
	edistance int
}

var testCases = []levTest{
	levTest{s1: "A", s2: "A", edistance: 0},
	levTest{s1: "A", s2: "B", edistance: 1},
	levTest{s1: "A", s2: "AB", edistance: 1},
	levTest{s1: "AA", s2: "AB", edistance: 1},
	levTest{s1: "AAA", s2: "AAB", edistance: 1},
	levTest{s1: "AAA", s2: "B", edistance: 3},
	levTest{s1: "AAAA", s2: "B", edistance: 4},
	levTest{s1: "kitten", s2: "sitting", edistance: 3},
	levTest{s1: "saturday", s2: "sunday", edistance: 3},
	levTest{s1: "exponential", s2: "polynomial", edistance: 6},
}

func TestDistance(t *testing.T) {
	for _, test := range testCases {
		if d := levenshteinDistance(test.s1, test.s2); d != test.edistance {
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

func TestDistanceIterative(t *testing.T) {
	for _, test := range testCases {
		if d := levenshteinIterative(test.s1, test.s2); d != test.edistance {
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
