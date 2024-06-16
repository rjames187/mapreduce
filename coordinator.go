package main

import (
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"path/filepath"
	"sync"
	"time"
)

type Job struct {
	Type     string
	FilePath string
	State    string
}

type Coordinator struct {
	JobQueue []*Job
	QueueLock *sync.Mutex
	TakenJobs map[string]*Job
}

func (c *Coordinator) coordinate(dir string, addr string) {
	// initialize a FIFO queue of Map jobs
	c.initMapJobs(dir)
	log.Printf("%d Map jobs waiting for workers", len(c.JobQueue))

	// listen for messages from workers
	rpc.Register(c)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listening error: ", err)
	}
	go http.Serve(l, nil)
	for {
		duration, _ := time.ParseDuration("1s")
		time.Sleep(duration)
	}
}

// RPC handler for a worker job request
func (c *Coordinator) RequestJob(args *RequestJobArgs, reply *RequestJobReply) error {
	log.Print("Received a job request")
	c.QueueLock.Lock()
	defer c.QueueLock.Unlock()
	// pop a job off the queue and send to the worker
	job := c.JobQueue[0]
	reply.Job = job
	c.JobQueue = c.JobQueue[1:]
	return nil
}

// create the initial map jobs by building a list of all input files
func (c *Coordinator) initMapJobs(dir string) {
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
	c.JobQueue = queue
}

// instantiate a new coordinator
func NewCoordinator() *Coordinator {
	return &Coordinator{
		TakenJobs: map[string]*Job{},
		QueueLock: &sync.Mutex{},
	}
}