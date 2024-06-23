package plugins

import (
	"crypto/rand"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Crash struct{}

func fault() {
	num, _ := rand.Int(rand.Reader, big.NewInt(1000))
	if num.Int64() < 330 {
		log.Fatal("Worker crashing ...")
	} else if num.Int64() < 660 {
		bigMS, _ := rand.Int(rand.Reader, big.NewInt(10000))
		ms := bigMS.Int64() + 5000
		log.Printf("Worker sleeping for %d ms", ms)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func (c Crash) Map(input string) []*KeyValue {
	fault()
	res := []*KeyValue{}
	for _, word := range strings.Fields(input) {
		pair := KeyValue{word, "1"}
		res = append(res, &pair)
	}
	return res
}

func (c Crash) Reduce(pairs []*KeyValue) []*KeyValue {
	fault()
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