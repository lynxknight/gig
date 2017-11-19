package distance

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

func distance(i, j int, s1, s2 *string, matrix map[mkey]int) int {
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
	a1 := distance(i, j-1, s1, s2, matrix) + 1
	// RDBG("MOVE LEFT")
	a2 := distance(i-1, j, s1, s2, matrix) + 1
	// RDBG("MOVE UP-LEFT")
	a3 := distance(i-1, j-1, s1, s2, matrix) + ne
	// RDBG("[considering those for min: %v; %v; %v;", a1, a2, a3)
	res := min(a1, a2, a3)
	// RDBG("RETURN RES for i=%v; j=%v; RES=%v", i, j, res)
	// RD -= 1
	matrix[key] = res
	return res
}

// Distance calculates Levenshtein distance between s1 and s2
func Distance(s1, s2 string) int {
	matrix := make(map[mkey]int)
	return distance(len(s1)-1, len(s2)-1, &s1, &s2, matrix)
}
