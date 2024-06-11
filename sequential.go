package main

import (
	"fmt"
	"io/fs"
	"log"
	"mapreduce/plugins"
	"os"
	"path/filepath"
)

func sequentialMapReduce(path string, plugin string) {
	functions := plugins.Plugins[plugin]
	intermediate_pairs := map[string][]string{}

	// read input files and pass into map function
	err := filepath.Walk(fmt.Sprintf("./%s", path), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		pairs := functions.Map(string(data))
		for _, pair := range pairs {
			key := pair[0]
			val := pair[1]
			group, ok := intermediate_pairs[key]
			if !ok {
				intermediate_pairs[key] = []string{val}
			} else {
				intermediate_pairs[key] = append(group, val)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// pass intermediate k/v pairs into reduce function
	for key, val := range intermediate_pairs {
		res := functions.Reduce(key, val)
		fmt.Printf("%s: %s", key, res)
	}
}