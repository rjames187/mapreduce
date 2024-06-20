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
	intermediatePairs := []*plugins.KeyValue{}

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
		intermediatePairs = append(intermediatePairs, functions.Map(string(data))...)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// pass intermediate k/v pairs into reduce function
	finalPairs := functions.Reduce(intermediatePairs)

	sort.Slice(finalPairs, func(i, j int) bool {
		return finalPairs[i].Key < finalPairs[j].Key
	})

	for _, p := range finalPairs {
		fmt.Printf("%v: %v\n", p.Key, p.Value)
	}
}