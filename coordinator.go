package main

import (
	"io/fs"
	"log"
	"path/filepath"
)

type Job struct {
	Type     string
	FilePath string
	State    string
}

func coordinate(dir string, addr string) {
	jobQueue := initMapJobs(dir)
	log.Print(jobQueue)
}

func initMapJobs(dir string) []*Job {
	queue := []*Job{}
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name()[0] != 'i' || info.IsDir() {
			return nil
		}
		job := &Job{
			Type: "map",
			FilePath: path,
			State: "queued",
		}
		queue = append(queue, job)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return queue
}