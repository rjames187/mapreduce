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
	counter := map[string]*KeyValue{}
	for _, p := range pairs {
		count, _ := strconv.Atoi(p.Value)
		_, present := counter[p.Key]
		if present {
			oldCount, _ := strconv.Atoi(counter[p.Key].Value)
			counter[p.Key].Value = strconv.Itoa(count + oldCount)
		} else {
			counter[p.Key] = &KeyValue{Key: p.Key, Value: strconv.Itoa(count)}
		}
	}
	res := []*KeyValue{}
	for _, kv := range counter {
		res = append(res, kv)
	}
	return res
}