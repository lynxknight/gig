package distance

import (
	"strings"
	"unicode/utf8"
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
		return base
	}
	diffsize := slen - sublen
	scores := make([]int, diffsize+1)
	for i := 0; i < diffsize-1; i++ {
		scores = append(scores, levenshteinDistance(s[i:i+sublen+1], substr))
	}
	return base + min(scores...)
}

func levenscore(str, target string) int {
	// We could have just computed levenshtein distance, but then "short"
	// strings would always win over long ones. For example, querying "data-v3"
	// would yield following result order:
	// 	  * develop (branch A, not preferrable)
	// 	  * feature/data_v3-memory-leak (branch B, preferrable)
	// We failed to find exact match on branch B, so we fallback on calculating
	// distance, but branch A would yield better score for almost any longer
	// substring due to its length (i.e. pure replacements will do). To fix this
	// we try to calculate "best distance between target and substrings of str".
	//
	// There are several ways to pick substrings:
	// 1) Sliding window of length len(target). Pick index X of str, and calculate
	// distance between str[X:X+len(target)] and target for each X that allows such
	// window size. Good method, but it lacks perfomance for long strings, so we
	// need to throw away "looks like bad" windows.
	//
	// 2) Search for all "almost occurences" of target[0] in str, let's say it is
	// str[X], and calculate distance between target and str[X:X+len(target)].
	// "Almost occurence" between X and Y is a case when either X == Y, or Y is any
	// of chars that are located near X on the QWERTY keyboard. Almost occurences
	// for char "h" are "tyugjbnm"

	if len(str) <= len(target) {
		return levenshteinIterative(str, target)
	}

	lenDiff := len(str) - len(target)
	targetStartByte := target[0]
	typosMaps := GetTyposMaps()
	scores := make([]int, 0)
	for bytePos := 0; bytePos < lenDiff; bytePos++ {
		if str[bytePos] == targetStartByte || typosMaps[targetStartByte][str[bytePos]] {
			window := str[bytePos : bytePos+len(target)]
			score := levenshteinIterative(window, target)
			// fmt.Println(window, score)
			scores = append(scores, score)
		}
	}
	if len(scores) == 0 {
		return levenshteinIterative(str, target)
	}
	return min(scores...)
}

func levenshteinIterative(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}

// GetScore finds out how much points does "substr" gain in "s", the lesser
// the better. Exact matches grant -10 points, if there are no "exact" matches,
// we try to go for "levenshtein" matches
func GetScore(s, substr string) int {
	if em := exactMatches(s, substr); em != 0 {
		return em * -10
	}
	score := 0
	if len(substr) > 3 {
		score = levenscore(s, substr)
	} else {
		score = levenshteinIterative(s, substr)
	}
	return score
}
