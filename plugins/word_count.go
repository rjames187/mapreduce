package plugins

import (
	"strconv"
	"strings"
)

type WordCount struct{}

func (wc WordCount) Map(input string) [][]string {
	res := [][]string{}
	for _, word := range strings.Fields(input) {
		pair := []string{word, "1"}
		res = append(res, pair)
	}
	return res
}

func (wc WordCount) Reduce(key string, values []string) []string {
	res := 0
	for _, v := range values {
		value, _ := strconv.Atoi(v)
		res += value
	}
	return []string{strconv.Itoa(res)}
}