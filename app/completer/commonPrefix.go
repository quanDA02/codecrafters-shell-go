package completer

// check common prefix
func CommonPrefix(line [][]rune) []rune {
	first, last := (line[0]), line[len(line)-1]
	result := first[:0]
	for i := 0; i < len(first) && i < len(last) && first[i] == last[i]; i++ {
		result = first[:i+1]
	}
	return result
}
