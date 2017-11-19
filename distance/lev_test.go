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

func TestDistanceEqualStringsLengthOne(t *testing.T) {
	if d := Distance("A", "A"); d != 0 {
		t.Errorf("'A' and 'A' should have distance of 0, got: %v", d)
	}
}

func TestDistanceLengthOneReplacement(t *testing.T) {
	if d := Distance("A", "B"); d != 1 {
		t.Errorf("'A' and 'B' should have distance of 1, got: %v", d)
	}
}

func TestDistanceAddition(t *testing.T) {
	if d := Distance("A", "AB"); d != 1 {
		t.Errorf("'A' and 'AB' should have distance of 1, got: %v", d)
	}
}

func TestDistanceAdditionLengthTwo(t *testing.T) {
	if d := Distance("AA", "AB"); d != 1 {
		t.Errorf("'A' and 'AB' should have distance of 1, got: %v", d)
	}
}

func TestDistanceReplacement(t *testing.T) {
	if d := Distance("AAA", "AAB"); d != 1 {
		t.Errorf("'A' and 'AB' should have distance of 1, got: %v", d)
	}
}

func Test12(t *testing.T) {
	if d := Distance("AAA", "B"); d != 3 {
		t.Errorf("'A' and 'AB' should have distance of 3, got: %v", d)
	}
}

func Test13(t *testing.T) {
	if d := Distance("AAAA", "B"); d != 4 {
		t.Errorf("'A' and 'AB' should have distance of 4, got: %v", d)
	}
}

func TestWiki(t *testing.T) {
	if d := Distance("EXPONENTIAL", "POLYNOMIAL"); d != 6 {
		t.Errorf("'A' and 'AB' should have distance of 6, got: %v", d)
	}
}

func TestWiki2(t *testing.T) {
	if d := Distance("kitten", "sitting"); d != 3 {
		t.Errorf("'A' and 'AB' should have distance of 3, got: %v", d)
	}
}

func TestWiki3(t *testing.T) {
	if d := Distance("saturday", "sunday"); d != 3 {
		t.Errorf("'A' and 'AB' should have distance of 3, got: %v", d)
	}
}
