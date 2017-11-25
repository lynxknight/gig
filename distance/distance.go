package distance

import (
	"fmt"
	"strings"
)

func min(ints ...int) int {
	m := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < m {
			m = ints[i]
		}
	}
	return m
}

// var RD = -1

// func RDBG(format string, a ...interface{}) {
// 	return
// 	for i := 0; i < RD; i++ {
// 		fmt.Printf("  ")
// 	}
// 	fmt.Printf(format, a...)
// 	fmt.Printf("\n")
// }

type mkey struct {
	a1, a2 int
}

func levenshteinDistanceRecursive(i, j int, s1, s2 *string, matrix map[mkey]int) int {
	// RD += 1
	key := mkey{i, j}
	v, ok := matrix[key]
	if ok {
		// RDBG("RETURN i=%v j=%v FROM CHACHE answer=%v", i, j, v)
		// RD -= 1
		return v
	}

	// RDBG("i=%v; j=%v;", i, j)
	// if RD > 300 {
	// os.Exit(1)
	// }
	if i == -1 && j == -1 {
		// RDBG("RETURN i=-1 j=-1 answer=0")
		// RD -= 1
		return 0
	}
	if i == -1 {
		// RDBG("RETURN i=-1 and answer j=%v", j)
		// RD -= 1
		return j + 1
	}
	if j == -1 {
		// RDBG("RETURN j=-1 and answer i=%v", i)
		// RD -= 1
		return i + 1
	}

	ne := 1
	if (*s1)[i] == (*s2)[j] {
		ne = 0
	}

	if i == 0 && j == 0 {
		// RDBG("RETURN ne: %v", ne)
		// RD -= 1
		return ne
	}
	// RDBG("MOVE UP")
	a1 := levenshteinDistanceRecursive(i, j-1, s1, s2, matrix) + 1
	// RDBG("MOVE LEFT")
	a2 := levenshteinDistanceRecursive(i-1, j, s1, s2, matrix) + 1
	// RDBG("MOVE UP-LEFT")
	a3 := levenshteinDistanceRecursive(i-1, j-1, s1, s2, matrix) + ne
	// RDBG("[considering those for min: %v; %v; %v;", a1, a2, a3)
	res := min(a1, a2, a3)
	// RDBG("RETURN RES for i=%v; j=%v; RES=%v", i, j, res)
	// RD -= 1
	matrix[key] = res
	return res
}

// LevenshteinDistance calculates Levenshtein distance between s1 and s2
func levenshteinDistance(s1, s2 string) int {
	matrix := make(map[mkey]int)
	return levenshteinDistanceRecursive(len(s1)-1, len(s2)-1, &s1, &s2, matrix)
}

func exactMatches(s, substr string) int {
	return strings.Count(s, substr)
}

func levenshteinScore(s, substr string) int {
	base := 100
	slen := len(s)
	sublen := len(substr)
	if slen < sublen {
		fmt.Println("slen", slen, "less than substr", sublen, "return base")
		return base
	}
	diffsize := slen - sublen
	scores := make([]int, diffsize+1)
	for i := 0; i < diffsize-1; i++ {
		fmt.Printf("  calculating substr for levenstien %v\n", s[i:i+sublen+1])
		scores = append(scores, levenshteinDistance(s[i:i+sublen+1], substr))
	}
	return base + min(scores...)
}

// GetScore finds out how much points does "substr" gain in "s", the lesser
// the better. Exact matches grant -10 points, if there are no "exact" matches,
// we try to go for "levenshtein" matches
func GetScore(s, substr string) int {
	if em := exactMatches(s, substr); em != 0 {
		fmt.Printf("Exact match found for {%v, %v} = %v\n", s, substr, em)
		return em * -10
	}
	score := levenshteinScore(s, substr)
	fmt.Printf("Calculating distance for {%v, %v} = %v\n", s, substr, score)
	return score
}
