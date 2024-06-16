package plugins

import (
	"strconv"
	"strings"
)

type WordCount struct{}

func (wc WordCount) Map(input string) []*KeyValue {
	res := []*KeyValue{}
	for _, word := range strings.Fields(input) {
		pair := KeyValue{word, "1"}
		res = append(res, &pair)
	}
	return res
}

func (wc WordCount) Reduce(pairs []*KeyValue) []*KeyValue {
	// all pairs in a reduce call have the same key
	key := pairs[0].Key
	res := 0
	for _, p := range pairs {
		count, _ := strconv.Atoi(p.Value)
		res += count
	}
	return []*KeyValue{&KeyValue{key, strconv.Itoa(res)}}
}