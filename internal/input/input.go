package input

import (
	"bufio"
	"strconv"
	"strings"
)

func ReadLine(r *bufio.Reader) string {
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func ReadInt(r *bufio.Reader) int {
	v, err := strconv.Atoi(ReadLine(r))
	if err != nil {
		return 1
	}
	return v
}

func ReadFloat(r *bufio.Reader) float64 {
	return ParseFloat(ReadLine(r))
}

func ParseFloat(s string) float64 {
	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0
	}
	return v
}
