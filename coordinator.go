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
	FilePaths []string
	Id int
	NReduce int
}

type Coordinator struct {
	JobQueue []*Job
	mu *sync.Mutex
	TakenJobs map[int]*Job
	ReduceFilePaths map[int][]string
	NReduce int
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

	// keep checking if all map jobs are done
	for {
		duration, _ := time.ParseDuration("1s")
		time.Sleep(duration)
		c.mu.Lock()
		if len(c.JobQueue) == 0 && len(c.TakenJobs) == 0 {
			log.Println("All map jobs completed ...")
			c.mu.Unlock()
			break
		}
		c.mu.Unlock()
	}
	
	// create FIFO queue of reduce jobs
	c.initReduceJobs()
	log.Printf("%d Reduce jobs waiting for workers", len(c.JobQueue))

	for {
		duration, _ := time.ParseDuration("1s")
		time.Sleep(duration)
		c.mu.Lock()
		if len(c.JobQueue) == 0 && len(c.TakenJobs) == 0 {
			log.Println("All reduce jobs completed ...")
			break
		}
		c.mu.Unlock()
	}
}

// RPC handler for a worker job request
func (c *Coordinator) RequestJob(args *RequestJobArgs, reply *RequestJobReply) error {
	log.Print("Received a job request")
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.JobQueue) == 0 {
		// if no jobs are queued, send signal to wait
		reply.Job = nil
		return nil
	}
	// pop a job off the queue and send to the worker
	job := c.JobQueue[0]
	reply.Job = job
	c.JobQueue = c.JobQueue[1:]
	c.TakenJobs[job.Id] = job
	go c.timeout(job)
	return nil
}

// RPC handler for complete map job 
func (c *Coordinator) CompleteMapJob(args *CompleteMapJobArgs, reply *CompleteMapJobReply) error {
	log.Println("Received complete map job call")
	c.mu.Lock()
	defer c.mu.Unlock()
	storedJob, ok := c.TakenJobs[args.Id]
	// handles the case where the job was already done by another worker
	if !ok || storedJob.Type == "reduce" {
		return nil
	}
	delete(c.TakenJobs, args.Id)
	for id, newPaths := range args.FilePaths {
		paths, present := c.ReduceFilePaths[id]
		if present {
			c.ReduceFilePaths[id] = append(paths, newPaths...)
		} else {
			c.ReduceFilePaths[id] = newPaths
		}
	}
	log.Println(c.ReduceFilePaths)
	return nil
}

// RPC handler for complete reduce job
func (c *Coordinator) CompleteReduceJob(args *CompleteReduceJobArgs, reply *CompleteReduceJobReply) error {
	log.Println("Received complete reduce job call")
	c.mu.Lock()
	defer c.mu.Unlock()
	storedJob, ok := c.TakenJobs[args.Id]
	// handles the case where the job was already done by another worker
	if !ok || storedJob.Type == "map" {
		return nil
	}
	delete(c.TakenJobs, args.Id)
	return nil
}

// procedure that executes upon a job timeout
func (c *Coordinator) timeout(job *Job) {
	time.Sleep(time.Duration(10) * time.Second)
	c.mu.Lock()
	defer c.mu.Unlock()
	storedJob, ok := c.TakenJobs[job.Id]
	if !ok {
		return
	}
	if storedJob.Type != job.Type {
		return
	}
	log.Printf("Timeout for %s job #%d ... assigning to new worker", job.Type, job.Id)
	delete(c.TakenJobs, job.Id)
	c.JobQueue = append(c.JobQueue, job)
}

// create the initial map jobs by building a list of all input files
func (c *Coordinator) initMapJobs(dir string) {
	queue := []*Job{}
	numJob := 1
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name()[0] != 'i' || info.IsDir() {
			return nil
		}
		job := &Job{
			Type: "map",
			FilePaths: []string{path},
			Id: numJob,
			NReduce: c.NReduce,
		}
		queue = append(queue, job)
		numJob += 1
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	c.JobQueue = queue
}

// create list of reduce jobs
func (c *Coordinator) initReduceJobs() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for j, paths := range c.ReduceFilePaths {
		job := &Job{
			Id: j,
			FilePaths: paths,
			Type: "reduce",
		}
		c.JobQueue = append(c.JobQueue, job)
	}
}

// instantiate a new coordinator
func NewCoordinator(nReduce int) *Coordinator {
	return &Coordinator{
		TakenJobs: map[int]*Job{},
		mu: &sync.Mutex{},
		ReduceFilePaths: map[int][]string{},
		NReduce: nReduce,
	}
}