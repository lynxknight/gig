package distance

import (
	"strings"
	"unicode/utf8"
)

func max2(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min2(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min(ints ...int) int {
	m := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < m {
			m = ints[i]
		}
	}
	return m
}

// Score struct contains actual score and indexes that represet this score
type Score struct {
	Distance, I1, I2 int // TODO: i1 i2 into range
}

const levenThreshold = 999999 // just random big value who cares

func exactMatches(s, substr string) Score {
	// TODO: rethink magic values
	index := strings.Index(s, substr)
	if index == -1 {
		return Score{Distance: -1}
	}
	// TODO: magic value for exact match was set to 0 since we have an allowed
	// range of minDistance(branches)+5 threshold (but maybe changing it was
	// not as good of an idea)
	return Score{0, index, index + len(substr)}
}

func calcLevenshteinScore(str, target string) Score {
	/* We could have just computed levenshtein distance, but then "short"
	strings would always win over long ones. For example, querying "data-v3"
	would yield following result order:
		* develop (branch A, not preferrable)
		* feature/data_v3-memory-leak (branch B, preferrable)
	We failed to find exact match for both branches, so we fallback on distance
	calculation. Problem is, branch A would yield better score for almost any
	longer substring due to its length (i.e. pure replacements will do). So,
	we are calculating "best distance between target and substrings of str".

	There are several ways to pick substrings:
	1) Sliding window of length len(target). Pick index X of str, and calculate
	distance between str[X:X+len(target)] and target for each X that allows such
	window size. Good method, but it lacks perfomance for long strings, so we
	need to throw away "looks like bad" windows.

	2) Search for all "almost occurences" of target[0] in str, let's say it is
	str[X], and calculate distance between target and str[X:X+len(target)].
	"Almost occurence" between char A and B is a case when either A == B, or
	B belongs to set of chars located near A on the QWERTY keyboard. For example,
	almost occurences for char "h" are "tyugjbnm" */
	var score Score
	score.Distance = levenThreshold
	if len(str) <= len(target) {
		distance := calcLeventshteinDistance(str, target, score.Distance)
		return Score{distance, 0, len(str)}
	}
	lenDiff := len(str) - len(target)
	targetStartByte := target[0]
	typosMaps := GetTyposMaps()
	for bytePos := 0; bytePos < lenDiff; bytePos++ {
		isAllowedTypo := typosMaps[targetStartByte][str[bytePos]]
		if str[bytePos] == targetStartByte || isAllowedTypo {
			window := str[bytePos : bytePos+len(target)]
			newDistance := calcLeventshteinDistance(window, target, score.Distance)
			if newDistance < score.Distance {
				score.Distance = newDistance
				score.I1, score.I2 = bytePos, bytePos+len(target)
			}
		}
	}
	if score.I2 != 0 { // basically we check if we hit "if" condition in "for" loop
		return score
	}
	return Score{calcLeventshteinDistance(str, target, score.Distance), 0, len(str)}
}

func calcLeventshteinDistance(a, b string, threshold int) int {
	// TODO: replace with naive implementation (or understand how this works)
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
// the better. Exact matches grant -10 points each, if there are no exact
// matches, we try to go for levenshtein distance
func GetScore(s, substr string) Score {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	// TODO: try using same terms everywhere. Current zoo of "a", "s", "str"
	// looks silly
	if em := exactMatches(s, substr); em.Distance != -1 {
		return em
	}
	var score Score
	if len(substr) > 3 {
		score = calcLevenshteinScore(s, substr)
	} else { // TODO why this "if" exists?
		score.Distance = calcLeventshteinDistance(s, substr, levenThreshold)
		score.I2 = min(len(substr), len(s))
	}
	return score
}
