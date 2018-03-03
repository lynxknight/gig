package distance

import (
	"strings"
)

type typoMap = map[byte]bool

var typoMapsGlobal map[byte]typoMap

// CreateTyposMaps creates typo maps
func CreateTyposMaps() map[byte]typoMap {
	// TODO: runes instead of bytes
	qwertyLayout := "1234567890-=\nqwertyuiop[]\nasdfghjkl;'\\\nzxcvbnm,./"
	lines := make(map[int]string)
	for index, line := range strings.Split(qwertyLayout, "\n") {
		lines[index] = line
	}
	typoMaps := make(map[byte]typoMap)
	for lineNumber, line := range lines {
		for bytePos := 0; bytePos < len(line); bytePos++ {
			typoMaps[line[bytePos]] = make(typoMap)
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if i == 0 && j == 0 {
						continue
					}
					// check bounds
					otherLine, ok := lines[lineNumber+i]
					if !ok || bytePos+j > len(otherLine)-1 || bytePos+j < 0 {
						continue
					}
					neighbourByte := otherLine[bytePos+j]
					typoMaps[line[bytePos]][neighbourByte] = true
				}
			}
		}
	}
	return typoMaps
}

// GetTyposMaps -
func GetTyposMaps() map[byte]typoMap {
	if typoMapsGlobal != nil {
		return typoMapsGlobal
	}
	typoMapsGlobal = CreateTyposMaps()
	return typoMapsGlobal
}

/*
from collections import defaultdict
nbgh = defaultdict(set)
for lineNum in mp:
    for charNum, char in enumerate(mp[lineNum]):
        for i in [-1, 0, 1]:
            for j in [-1, 0, 1]:
                if i == 0 and j == 0:
                    continue
                try:
                    nbgh[char].add(mp[lineNum+i][charNum+j])
                except (KeyError, IndexError):
					pass
*/
