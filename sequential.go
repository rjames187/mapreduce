package main

import (
	"fmt"
	"io/fs"
	"log"
	"mapreduce/plugins"
	"os"
	"path/filepath"
	"sort"
)

func sequentialMapReduce(path string, plugin string) {
	functions := plugins.Plugins[plugin]
	intermediate_pairs := map[string][]*plugins.KeyValue{}

	// read input files and pass into map function
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name()[0] != 'i' {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		pairs := functions.Map(string(data))
		for _, pair := range pairs {
			key := pair.Key
			val := pair.Value
			group, ok := intermediate_pairs[key]
			if !ok {
				intermediate_pairs[key] = []*plugins.KeyValue{{Key: key, Value: val,}}
			} else {
				intermediate_pairs[key] = append(group, &plugins.KeyValue{Key: key, Value: val,})
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// pass intermediate k/v pairs into reduce function
	final_pairs := []*plugins.KeyValue{}
	for _, group := range intermediate_pairs {
		res := functions.Reduce(group)
		final_pairs = append(final_pairs, res...)
	}

	sort.Slice(final_pairs, func(i, j int) bool {
		return final_pairs[i].Key < final_pairs[j].Key
	})

	for _, p := range final_pairs {
		fmt.Printf("%v: %v\n", p.Key, p.Value)
	}
}