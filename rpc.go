package main

type RequestJobArgs struct{}

type RequestJobReply struct {
	Job *Job
}

type CompleteMapJobArgs struct {
	Id        int
	FilePaths map[int][]string
}

type CompleteMapJobReply struct{}